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
            valid_attributes (string[]): all valid attributes. The issuer
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
        # sk = [x, y1, y2, y3, ...]
        # pk = [X, Y1, Y2, Y3, ...]
        for _ in range(len(valid_attributes) +1):                     # +1 because of X        
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
        
        return bytearray(jsonpickle.encode({"g": self.g, "gt": self.gt, "G1": self.G1, "G2": self.G2,
                                            "GT": self.GT, "valid_attr": self.valid_attributes,
                                            "pk": self.pk[1:], "pkt": self.pkt}), 'utf-8')
        
        
    def get_serialized_secret_key(self):
        """Returns the secret key of the issuer.

        Args:
            No input

        Returns:
            byte[]: issuer's secret params and key
        """

        return bytearray(jsonpickle.encode({"sk": self.sk, "X": self.pk[0], "public_params": {"g": self.g, 
                                            "gt": self.gt, "G1": self.G1, "G2": self.G2, "GT": self.GT,
                                            "valid_attr": self.valid_attributes, "pk": self.pk[1:],
                                            "pkt": self.pkt}}), 'utf-8')


    @staticmethod
    def issue(C, user, revealed_attr, server_sk):
        """Issues a credential for a new user. 

        This function should receive a issuance request from the user
        (AnonCredential.create_issue_request), and a list of known attributes of the
        user (e.g. the server received bank notes for subscriptions x, y, and z).

        Args:
            C (G1Element): commit from user
            user (string): username of the user
            revealed_attr (private_attr): list of revealed attributes of the user
            server_sk (dict): the server's secret key and other parameters

        Returns:
            sigma (tuple(G1Element, G1Element)): credentials of the issuer
        """

        server_pb_params = server_sk.get("public_params")

        # (u, X)
        u = petrelic.bn.Bn.from_num(rd.randint(1, server_pb_params.get("G1").order()))
        
        X = server_sk.get("X")
        Y = server_pb_params.get("pk")
        # create credentials sigma'
        sigma2 = X * C
        for i in range(len(revealed_attr)):
            sigma2 *= (Y[i] ** int(revealed_attr[i]))

        sigma = [server_pb_params.get("g") ** u, sigma2 ** u]

        print("Credential issuance done")
        return sigma


class AnonCredential(object):
    """An AnonCredential"""

    def create_issue_request(self, private_attributes, G1, Y, g, t, pub_attr_len):
        """Gets all known attributes (subscription) of a user and creates an issuance request.
        You are allowed to add extra attributes to the issuance.

        Args:
            private_attributes (string[]): list of private attributes of the client
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
            # Y[i] from 0 to pub_attr_len are pk for public attributes
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

        r0 = (v0 + c_hash*t) % G1.order()
        
        r = [r0]
        for vi, a in zip(v[1:], private_attributes):
            ri = (vi + c_hash*int(a)) % G1.order()
            r.append(ri)
        
        print("Issue request created")
        return C, gamma, r


    def receive_issue_response(self, sigma, t):
        """This function finishes the credential based on the response of issue.

        Hint: you need both secret values from the create_issue_request and response
        from issue to build the credential.

        Args:
            sigma (G1Element[]): response sigma' from issuer
            t (G1Element): secret t used for user commit

        Returns:
            G1Element[]: finalized credentials
        """

        return [sigma[0], sigma[1]*(sigma[0]**(-t))]



class Signature(object):
    """A Signature"""

    def create_sign_request(self, server_pk, sigma, message, revealed_attr, secret_attr):
        """Signs the message.

        Args:
            server_pk (dict): the server's public key and other parameters
            revealed_attr (string []): a list of revealed attributes
            sigma (G1Element[]): client's credential
            message (byte[]): message to sign
            revealed_attr (string[]): list of revealed attributes
            secret_attr (string[]): list of secret attributes
        Returns:
            t (Bn): secret parameter used to create signature
        """

        G1 = server_pk.get("G1")
        gt = server_pk.get("G2").generator()
        Yt = server_pk.get("pkt")
        Yt = Yt[1:]
        r, t = petrelic.bn.Bn.from_num(rd.randint(1, G1.order())), petrelic.bn.Bn.from_num(rd.randint(1, G1.order()))
        o = [sigma[0] ** r, (sigma[1]*(sigma[0]**t))**r]

        self.message = message.decode()

        # C = e(o1, gt^t) PROD in H (e(o1, Yt)^ai)
        C = o[0].pair(gt**t) 
        
        v = [petrelic.bn.Bn.from_num(rd.randint(1, G1.order()))]
        gamma = [o[0].pair(gt)**v[0]]

        for i in range(len(secret_attr)):
            C *= (o[0].pair(Yt[i+ len(revealed_attr)]))**(petrelic.bn.Bn.from_num(secret_attr[i]))
            vi = petrelic.bn.Bn.from_num(rd.randint(1, G1.order()))
            v.append(vi)
            gamma.append(o[0].pair(Yt[i + len(revealed_attr)])**vi)

        to_hash = message + C.to_binary() + gt.to_binary()
        for i in range(len(Yt)):            
            to_hash += Yt[i].to_binary()

        for gammai in gamma:            
            to_hash += gammai.to_binary()

        # c = H(message, C, gt, Yt1, Yt2, ..., gamma0, gamma1, gamma2)
        c_hash = int(hashlib.sha512(to_hash).hexdigest(), 16)
        
        r0 = (v[0] + c_hash*t) % G1.order()        
        r = [r0]  
        for vi, a in zip(v[1:], secret_attr):
            ri = (vi + c_hash*int(a)) % G1.order()         
            r.append(ri)       
        
        self.sigma = o
        self.C = C
        self.gamma = gamma
        self.r = r
        print("Sign request created")

        return t
        
        

    def verify(self, issuer_public_info, public_attrs, message):
        """Verifies a signature.

        Args:
            issuer_public_info (): output of issuer's 'get_serialized_public_key' method
            public_attrs (string[]): public attributes
            message (byte []): list of messages

        returns:
            valid (boolean): is signature valid
        """

        G1 = issuer_public_info.get("G1")
        GT = issuer_public_info.get("GT")
        gt = issuer_public_info.get("G2").generator()
        Yt = issuer_public_info.get("pkt")
        Xt = Yt[0]
        Yt = Yt[1:]

        if self.sigma[0] == G1.neutral_element():
            print("Sigma0 is neutral element")
            return False
        
        to_hash = message + self.C.to_binary() + gt.to_binary()
        
        left = self.sigma[1].pair(gt)
        left *= (self.sigma[0].pair(Xt))**(-1)
        for i in range(len(public_attrs)):
            a = (self.sigma[0].pair(Yt[i])**(-int(public_attrs[i])))
            left *= a

        # print("!!!!!!!!!!!!!!!!!!!")
        # print("equality test holds", left == self.C)
        # print("!!!!!!!!!!!!!!!!!!!")

        for i in range(len(Yt)):            
            to_hash += Yt[i].to_binary()

        gamma_prod = GT.neutral_element()
        for gammai in self.gamma:
            to_hash += gammai.to_binary()
            gamma_prod *= gammai

        # c = H(message, C, gt, Yt1, Yt2, ..., gamma0, gamma1, gamma2)
        c_hash = int(hashlib.sha512(to_hash).hexdigest(), 16)        

        #forget left and do as with C
        check_prod = (left**((-c_hash) % G1.order()))*(self.sigma[0].pair(gt)**self.r[0])
        for i in range(len(self.r)-1):
            check_prod *= self.sigma[0].pair(Yt[i + len(public_attrs)])**self.r[i+1]
        
        print("ZKP is valid ", check_prod == gamma_prod)
        return check_prod == gamma_prod
        

    def serialize(self):
        """Serialize the object to a byte array.

        Returns: 
            byte[]: a byte array 
        """
        
        return jsonpickle.encode(self).encode('utf-8')

    @staticmethod
    def deserialize(data):
        """Deserializes the object from a byte array.

        Args: 
            data (byte[]): a byte array 

        Returns:
            Signature
        """

        if type(data) == str:
            return jsonpickle.decode(data)
        else:
            return jsonpickle.decode(data.decode('utf-8'))




