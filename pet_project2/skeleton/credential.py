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
        for i in range(len(valid_attributes) + 2):
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

        You should design the issue_request as you see fit.
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

        You should design the issue_request as you see fit.
        """
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

    
    def create_sign_request(self, server_pk, sigma, message, revealed_attr, secret_attr):

        G1 = server_pk.get("G1")
        GT = server_pk.get("GT")
        gt = server_pk.get("G2").generator()
        Yt = server_pk.get("pkt")
        Xt = Yt[0]
        Yt = Yt[1:]
        r, t = petrelic.bn.Bn.from_num(rd.randint(1, G1.order())), petrelic.bn.Bn.from_num(rd.randint(1, G1.order()))
        o = [sigma[0] ** r, (sigma[1]*(sigma[0]**t)**r)]
        
        self.message = message.decode()
        
        C = o[1].pair(gt**t) 
        
        v = [petrelic.bn.Bn.from_num(rd.randint(1, GT.order()))]
        gamma = [o[0].pair(gt)**v[0]]
        for i in range(len(secret_attr)):
            C *= (o[0].pair(Yt[i+ len(revealed_attr)]))**(int(secret_attr[i]))
            vi = petrelic.bn.Bn.from_num(rd.randint(1, GT.order()))
            v.append(vi)
            gamma.append(o[0].pair(Yt[i + len(revealed_attr)])**vi)
        to_hash = message + C.to_binary() + gt.to_binary()
        for i in range(len(revealed_attr)):            
            to_hash += Yt[i  + len(revealed_attr)].to_binary()

        for gammai in gamma:            
            to_hash += gammai.to_binary()
        # c = H(message, C, gt, Yt1, Yt2, ..., gamma0, gamma1, gamma2)
        c_hash = int(hashlib.sha512(to_hash).hexdigest(), 16)
        print(c_hash)
        r0 = v[0] - c_hash*t
        if r0 < 0:
            r0 = r0 % GT.order()
        r = [r0]  

        for vi, a in zip(v[1:], secret_attr):
            res = (vi - c_hash*int(a))
            
            if  res < 0:
                res = res % GT.order()
                
            r.append(res)
       
        self.sigma = o
        '''
        C = o[0].pair(gt)**t
        C *= o[0].pair(Xt)
        attr = revealed_attr + secret_attr
        for i in range(len(attr)):
            C *= o[0].pair(Yt[i])** int(attr[i])
        '''
        print("sign request created")
        self.C = C
        self.gamma = gamma
        self.r = r
        
        

    def verify(self, issuer_public_info, public_attrs, message):
        """Verifies a signature.

        Args:
            issuer_public_info (): output of issuer's 'get_serialized_public_key' method
            public_attrs (dict): public attributes
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
            print("sigma0 is neutral element")
            return False
        
        to_hash = message + self.C.to_binary() + gt.to_binary()
        

        print("length should be 2 ", len(self.r))
        for i in range(len(self.r)):            
            to_hash += Yt[i + len(public_attrs)].to_binary()

        gamma_prod = GT.neutral_element()
        for gammai in self.gamma:
            to_hash += gammai.to_binary()
            gamma_prod *= gammai

        
        # c = H(message, C, gt, Yt1, Yt2, ..., gamma0, gamma1, gamma2)
        
        c_hash = int(hashlib.sha512(to_hash).hexdigest(), 16)
        print(c_hash)
        print('-----------------------------')


        #check left of inequality 
        left = self.sigma[1].pair(gt)* (self.sigma[0].pair(Xt)**(-1))
        for i in range(len(public_attrs)):
            left *= self.sigma[0].pair(Yt[i])**(- int(public_attrs[i]))

        print("left is same as received C ", left == self.C)


        #forget left and do as with C
        check_prod = (self.C**c_hash)*(self.sigma[0].pair(gt)**self.r[0])
        for i in range(len(self.r)-1):
            check_prod *= self.sigma[0].pair(Yt[i + len(public_attrs)])**self.r[i+1]
        # zero knowledge proof with: proof, sigmaPrime,...
        print(check_prod)
        print("-------------------")
        print(gamma_prod)
        print("------------------------")
        print(left)
        print("gamma prod is same as check prod",  gamma_prod == check_prod )
        return None
        '''

        left = self.sigma[1].pair(gt)
        left *= self.sigma[0].pair(Xt)**(-1)
        for i in range(len(public_attrs)):
            left *= self.sigma[0].pair(Yt[i])**(-int(public_attrs[i]))

        print(self.sigma[1].pair(gt))
        print('-------------------')
        print(self.C)

        return (self.sigma[1].pair(gt)) == self.C
        '''
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




