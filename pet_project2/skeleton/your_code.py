"""
Classes that you need to complete.
"""

# Optional import
from serialization import jsonpickle
import credential
import petrelic
import random as rd
class Server:
    """Server"""
    

    @staticmethod
    def generate_ca(valid_attributes):
        """Initializes the credential system. Runs exactly once in the
        beginning. Decides on schemes public parameters and choses a secret key
        for the server.

        Args:
            valid_attributes (string): a list of all valid attributes. Users cannot
            get a credential with a attribute which is not included here.

            Note: You can use JSON to encode valid_attributes in the string.

        Returns:
            (tuple): tuple containing:
                byte[] : server's pubic information
                byte[] : server's secret key
            You are free to design this as you see fit, but all commuincations
            needs to be encoded as byte arrays.
        """
        Issuer = credential.Issuer()
        Issuer.setup(valid_attributes.split(","))
        pub = Issuer.get_serialized_public_key()
        sec = Issuer.get_serialized_secret_key()

        return pub, sec

    @staticmethod
    def register(server_sk, issuance_request, username, attributes):
        """ Registers a new account on the server.

        Args:
            server_sk (byte []): the server's secret key (serialized)
            issuance_request (bytes[]): The issuance request (serialized)
            username (string): username
            attributes (string): attributes

            Note: You can use JSON to encode attributes in the string.

        Return:
            response (bytes[]): the client should be able to build a credential
            with this response.
        """
        server_sk = jsonpickle.decode(server_sk.decode("utf-8"))
        #server_pk = server_sk.get("public_params")
        print("issuance req is ", issuance_request, type(issuance_request))
        issuance_request = jsonpickle.decode(issuance_request)
        sigma = credential.Issuer.issue(issuance_request.get("C"), username, "XX", server_sk)
        print('hi', sigma)
        
        return bytearray(jsonpickle.encode(sigma), 'utf-8')
    def check_request_signature(
        self, server_pk, message, revealed_attributes, signature
    ):
        """

        Args:
            server_pk (byte[]): the server's public key (serialized)
            message (byte[]): The message to sign
            revealed_attributes (string): revealed attributes
            signature (bytes[]): user's autorization (serialized)

            Note: You can use JSON to encode revealed_attributes in the string.

        Returns:
            valid (boolean): is signature valid
        """
        raise NotImplementedError


class Client:
    """Client"""

    def prepare_registration(self, server_pk, username, attributes):
        """Prepare a request to register a new account on the server.

        Args:
            server_pk (byte[]): a server's public key (serialized)
            username (string): username
            attributes (string): user's attributes

            Note: You can use JSON to encode attributes in the string.

        Return:
            tuple:
                byte[]: an issuance request
                (private_state): You can use state to store and transfer information
                from prepare_registration to proceed_registration_response.
                You need to design the state yourself.
        """
        server_pk = jsonpickle.decode(server_pk.decode("utf-8"))

        g = server_pk.get("g")
        Yi = server_pk.get("pk")
        G1= server_pk.get("G1")
        t = petrelic.bn.Bn.from_num(rd.randint(1, G1.order()))
        sk = petrelic.bn.Bn.from_num(rd.randint(1, G1.order()))
        attributes_list = attributes.split(",")
        attributes_list.append(sk)
        
        AnonCredential = credential.AnonCredential()
        C, gamma, r = AnonCredential.create_issue_request(attributes_list, G1, Yi, g, t)

        request = (jsonpickle.encode({"C": C, "gamma": gamma, "r": r})).encode('utf-8')
        private_state = {"C": C, "t": t}
        return (request, private_state) 

    def proceed_registration_response(self, server_pk, server_response, private_state):
        """Process the response from the server.

        Args:
            server_pk (byte[]): a server's public key (serialized)
            server_response (byte[]): the response from the server (serialized)
            private_state (private_state): state from the prepare_registration
            request corresponding to this response

        Return:
            credential (byte []): create an attribute-based credential for the user
        """
        sigma = jsonpickle.decode(server_response.decode('utf-8'))
        sigma = [sigma[0], sigma[1]/(sigma[0]**private_state.get("t"))]
        
        return bytearray(jsonpickle.encode(sigma), 'utf-8')

    def sign_request(self, server_pk, credential, message, revealed_info):
        """Signs the request with the clients credential.

        Arg:
            server_pk (byte[]): a server's public key (serialized)
            credential (byte[]): client's credential (serialized)
            message (byte[]): message to sign
            revealed_info (string): attributes which need to be authorized

            Note: You can use JSON to encode revealed_info.

        Returns:
            byte []: message's signature (serialized)
        """
        raise NotImplementedError
