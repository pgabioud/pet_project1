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
import hashlib
from petrelic.multiplicative.pairing import G1, G2, GT
class PSSignature(object):
    """PS's Multi-message signature from section 4.2
    
    **Important** This class has no direct use in the project.

    Implementing this class allows you to get familiar with coding crypto schemes
    and its simplicity in comparison with the ABC scheme allows you to realize
    misunderstandings/problems early on.
    """

    #Assume:
    L = 2

    def generate_key(self):
        g = G2.generator()
        sk = []
        pk = []
        for _ in range(self.L + 1):
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
        # (p, g, g tilde, G1, G2, GT, secret key, public key tilde, public key)
        self.g = G1.generator()
        self.gt = G2.generator()
        self.G1 = G1
        self.G2 = G2
        self.GT = GT
        self.valid_attributes = valid_attributes
        self.sk = []
        self.pkt = []
        self.pk = []
        for _ in range(len(valid_attributes) + 1):
            y = petrelic.bn.Bn.from_num(rd.randint(1, self.G2.order()))
            self.sk.append(y)
            self.pk.append(self.g**y)
            self.pkt.append(self.gt**y)

        print("Setup done")


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

        Args:
            C (G1Element): commit from user
            user (string): username of the user
            revealed_attr (list): list of revealed attributes of the user
            server_sk (dict): the server's secret key and other parameters

        Returns:
            sigma (tuple(G1Element, G1Element)): credentials of the issuer
        """

        rd.seed(user)
        server_pb_params = server_sk.get("public_params")

        # (u, X)
        u = petrelic.bn.Bn.from_num(rd.randint(1, server_pb_params.get("G1").order()))
        X = server_sk.get("X")

        # create credentials sigma'
        sigma2 = X * C
        for i in range(len(revealed_attr)):
            Y = server_pb_params.get("pk")[i]
            sigma2 *= Y ** petrelic.bn.Bn.from_num(int(revealed_attr[i]))

        sigma = [server_pb_params.get("g") ** u, sigma2 ** u]

        print("Credentials created")
        return sigma


class AnonCredential(object):
    """An AnonCredential"""

    def create_issue_request(self, private_attributes, G1, Y, g, t, pub_attr_len):
        """Gets all known attributes (subscription) of a user and creates an issuance request.
        You are allowed to add extra attributes to the issuance.

        Args:
            private_attributes (list): list of private attributes of the client
            G1 (petrelic group): group of the issueer private key
            Y (G1Element[]): public key of the issuer
            g (G1Element): generator of G1
            t (Bn): random big number between 1 and order of G1
            pub_attr_len (int): number of public attributes of the client

        Returns:
            C (G1Element): commit of the client
            gamma (G1Element[]): list of gamma used for zero knowledge proof of commit
            r (Bn): list of exponent for zero knowledge proof of commit
        """

        # Commit is C = g^t * Y1^a1 * Y2*a2 * ...

        # Zero Knowledge Proof:
        # - pick v0 on G1, compute gamma0 = g^v0
        # - pick vi on G1, compute gammai = Yi^vi for i from 1 to nb of private_attr
        # - compute c = Hash(C, g, Y1, Y2, ..., gamma0, gamma1, gamma2)
        # - compute r0 = v0 - c*t
        # - compute ri = vi - c*ai for i from 1 to nb of private_attr
        # => check: gamma0 * gamma1 * gamma2 * ... = C^c * g^r0 * Y1^r1 * Y2^r2 * ...

        C = g**t
        v = []
        gamma = []
        
        v0 = petrelic.bn.Bn.from_num(rd.randint(1, G1.order()))
        v.append(v0)
        gamma.append(g**v0)
        for i in range(len(private_attributes)):
            C *= Y[i+pub_attr_len]**int(private_attributes[i])
            vi = petrelic.bn.Bn.from_num(rd.randint(1, G1.order()))
            v.append(vi)
            gamma.append(Y[i+pub_attr_len]**vi)
        
        to_hash = C.to_binary() + g.to_binary()
        for i in range(len(private_attributes)):
            to_hash += Y[i+pub_attr_len].to_binary()
        for gammai in gamma:
            to_hash += gammai.to_binary()

        c_hash = int(hashlib.sha512(to_hash).hexdigest(), 16)
        
        r0 = v0 - c_hash*t
        if r0 < 0:
            r0 = r0 % G1.order()
        r = [r0]
        for vi, a in zip(v[1:], private_attributes):
            ri = (vi - c_hash*int(a))
            if  ri < 0:
                ri = ri % G1.order()
            r.append(ri)

        print("Issue request created")
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




