# The goal of this skeleton is helping you start with the credential.
# Following this API is not mandatory and you can change it as you see fit.
# This skeleton only provides major classes/functionality that you need. You 
# should define more classes and functions.

# Hint: for a clean code, you should create classes for messages that you want
# to pass between user and issuer. The serialization helps you with (de)serializing them
# (network API expects byte[] as input).

from serialization import jsonpickle
import random as rd
import petrelic
from petrelic.multiplicative.pairing import G1, G2
class PSSignature(object):

    L = 2
    """PS's Multi-message signature from section 4.2
    
    **Important** This class has no direct use in the project.

    Implementing this class allows you to get familiar with coding crypto schemes
    and its simplicity in comparison with the ABC scheme allows you to realize
    misunderstandings/problems early on.
    """

    def generate_key(self):
        g = G2.generator()
        sk = []
        pk = []
        for i in range(self.L + 1):
            y = petrelic.bn.Bn.from_num(rd.randint(1, G2.order()))
            sk.append(y)
            pk.append(g**y)
            
        return sk, pk

    def sign(self, sk, messages):             
        exp = G1.order().random()
        h = G1.generator() ** exp
        sum = sk[0]
        for i in range(len(sk)-1):
            rd.seed(messages[i])
            sum += sk[i+1] * rd.randint(1, G1.order())
        return [h, h**sum]


    def verify(self, pk, signature, messages):
        if signature[0] == G1.neutral_element():
            return False

        prod = pk[0]
        for i in range(len(pk)-1):
            rd.seed(messages[i])
            prod *= pk[i+1]** rd.randint(1, G1.order())
        if signature[0].pair(prod) != signature[1].pair(G2.generator()):
            return False
        else:
            return True


class Issuer(object):
    """Allows the server to issue credentials"""

    def setup(self, valid_attributes):
        """Decides the public parameters of the scheme and generates a key for
        the issuer.

        Args:
            valid_attributes (string): all valid attributes. The issuer
            will never be called with a value outside this list
        """
        # (p, g, G1, G2, GT, our u)
        self.g = G1.generator()
        self.gt = G2.generator()
        self.G1 = G1
        self.G2 = G2
        self.GT = petrelic.multiplicative.pairing.GT
        self.valid_attributes = valid_attributes
        self.sk = []
        self.pkt = []
        self.pk = []
        for i in range(len(valid_attributes) + 1):
            y = petrelic.bn.Bn.from_num(rd.randint(1, self.G2.order()))
            self.sk.append(y)
            self.pk.append(self.g**y)
            self.pkt.append(self.gt**y)


    def get_serialized_public_key(self):
        """Returns the public parameters and the public key of the issuer.

        Args:
            No input

        Returns:
            byte[]: issuer's public params and key
        """
        
        return bytearray(jsonpickle.encode({"g": self.g, "G1": self.G1, "G2": self.G2, "GT": self.GT, 
                                            "valid_attr": self.valid_attributes, "pk": self.pk[1:], 
                                            "pkt": self.pkt}), 'utf-8')
        
        
    def get_serialized_secret_key(self):
        """Returns the secret key of the issuer.

        Args:
            No input

        Returns:
            byte[]: issuer's secret params and key
        """
        return bytearray(jsonpickle.encode({"sk": self.sk, "X": self.pk[0], "public_params": {"g": self.g, "G1": self.G1, 
                                            "G2": self.G2, "GT": self.GT, "valid_attr": self.valid_attributes, 
                                            "pk": self.pk[1:], "pkt": self.pkt}}), 'utf-8')

    @staticmethod
    def issue(C, user, revealed_attr, server_sk):
        """Issues a credential for a new user. 

        This function should receive a issuance request from the user
        (AnonCredential.create_issue_request), and a list of known attributes of the
        user (e.g. the server received bank notes for subscriptions x, y, and z).

        You should design the issue_request as you see fit.
        """
        rd.seed(user)
        u = petrelic.bn.Bn.from_num(rd.randint(1, server_sk.get("public_params").get("G1").order()))
        X = server_sk.get("X")
        sigma2 = X * C
        for i in range(len(server_sk.get("public_params").get("valid_attr"))):
            if revealed_attr[i] != 'X':
                Y = server_sk.get("public_params").get("pk")[i]
                sigma2 *= Y ** petrelic.bn.Bn.from_num(revealed_attr[i])

        sigma = [server_sk.get("public_params").get("g") ** u, sigma2 ** u]
        return sigma


class AnonCredential(object):
    """An AnonCredential"""

    def create_issue_request(self, attributes, G1, Y, g, t):
        """Gets all known attributes (subscription) of a user and creates an issuance request.
        You are allowed to add extra attributes to the issuance.

        You should design the issue_request as you see fit.
        """
        C = g**t
        v = []
        gamma = []
        v0 = petrelic.bn.Bn.from_num(rd.randint(1, G1.order()))
        v.append(v0)
        gamma.append(g**v0)
        for i in range(len(attributes)):
            C *= Y[i]**int(attributes[i])
            vi = petrelic.bn.Bn.from_num(rd.randint(1, G1.order()))
            v.append(vi)
            gamma.append(Y[i]**vi)
        
        to_hash = []
        to_hash.append(C).append(g)
        for Yi in Y:
            to_hash.append(Yi)
        for gammai in gamma:
            to_hash.append(gammai)

        # c = H(C, g, Y1, Y2, ..., gamma0, gamma1, gamma2)
        c_hash = hash(tuple(to_hash))

        r = []
        for vi, a in zip(v, attributes):
            r.append(vi - c_hash*a)

        return C, gamma, r

        

    def receive_issue_response(self, sigma = [0,0], t= 1):
        """This function finishes the credential based on the response of issue.

        Hint: you need both secret values from the create_issue_request and response
        from issue to build the credential.

        You should design the issue_request as you see fit.
        """
        
        return Signature(sigma[0], sigma[1]/ sigma[0]**t)

    def sign(self, message, revealed_attr, sigma):
        """Signs the message.

        Args:
            message (byte []): message
            revealed_attr (string []): a list of revealed attributes

        Return:
            Signature: signature
        """
        toHash = message.decode("utf-8")
        toHash += str(sigma[0]) + str(sigma[1])
        for attr in revealed_attr:
            toHash += attr
        hashSign = hash(toHash)
        return hashSign


class Signature(object):
    """A Signature"""

    def __init__(self, sigma0, sigma1):
        self.sigma = [sigma0, sigma1]

    def verify(self, issuer_public_info, public_attrs, message, sigmaPrime, hashSignMsg, proof):
        """Verifies a signature.

        Args:
            issuer_public_info (): output of issuer's 'get_serialized_public_key' method
            public_attrs (dict): public attributes
            message (byte []): list of messages

        returns:
            valid (boolean): is signature valid
        """

        toHash = message.decode("utf-8")
        toHash += str(sigmaPrime[0]) + str(sigmaPrime[1])
        for attr in public_attrs:
            toHash += attr
        hashSign = hash(toHash)
        if hashSignMsg != hashSign:
            return False

        # zero knowledge proof with: proof, sigmaPrime,...


    def serialize(self):
        """Serialize the object to a byte array.

        Returns: 
            byte[]: a byte array 
        """
        return bytearray(jsonpickle.encode(self.sigma), "utf-8")

    @staticmethod
    def deserialize(data):
        """Deserializes the object from a byte array.

        Args: 
            data (byte[]): a byte array 

        Returns:
            Signature
        """
        return Signature(jsonpickle.decode(data.decode("utf-8"))[0], jsonpickle.decode(data.decode("utf-8"))[1])




