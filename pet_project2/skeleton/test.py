import pytest
from credential import PSSignature, Issuer, AnonCredential, Signature
from serialization import jsonpickle
import random as rd
import petrelic, hashlib
from petrelic.multiplicative.pairing import G1
from your_code import Client, Server


@pytest.fixture
def input_issuer():
    issuer = Issuer()
    issuer.setup(["male", "legal", "private_attr"])
    return issuer

@pytest.fixture
def input_anoncred():
    anoncred = AnonCredential()
    return anoncred

@pytest.fixture
def input_server():
    server = Server()
    return server

@pytest.fixture
def input_client():
    client = Client()
    return client

@pytest.fixture
def input_cred_params():
    t = petrelic.bn.Bn.from_num(rd.randint(1, G1.order()))
    private_attr = [petrelic.bn.Bn.from_num(rd.randint(1, G1.order()))]
    revealed_attr = [1,1]
    return t, private_attr, revealed_attr

#@pytest.mark.skip
def test_pssignature():
    machin = PSSignature()
    sk, pk = machin.generate_key()
    print("secret is ", sk[0], type(sk[0]))
    print("---------------------------------")
    signature = machin.sign(sk, [b"Fake news", b"apchii"])
    print('signature is ', signature)
    print("---------------------------------")
    assert machin.verify(pk, signature, [b"Fake news", b"apchii"])

#@pytest.mark.skip
def test_issue_request(input_server, input_issuer, input_anoncred, input_cred_params):
    valid_attr_str = "male, legal, sk"
    revealed_attr_str = "1,1"
    server_pk, _ = input_server.generate_ca(valid_attr_str)
    attributes_list = revealed_attr_str.split(",")
    server_pk = Signature.deserialize(server_pk)
    g = server_pk.get("g")
    Y = server_pk.get("pk")
    G1 = server_pk.get("G1")
    t = petrelic.bn.Bn.from_num(rd.randint(1, G1.order()))
    private_attributes = [petrelic.bn.Bn.from_num(rd.randint(1, G1.order()))]
        
    anonCredential = AnonCredential()
    #C, gamma, r = anonCredential.create_issue_request(private_attributes, G1, Y, g, t, len(attributes_list))
    
    #G1 = input_issuer.G1
    g1 = G1.generator()
    g1 = g
    print(g1 == g)
    """
    t = input_cred_params[0]
    pub_attr_len = 2
    private_attributes = input_cred_params[1]
    Y = input_issuer.pk[1:]
    """
    C, gamma, r = input_anoncred.create_issue_request(private_attributes, G1, Y, g1, t, len(attributes_list))

    to_hash = C.to_binary() + g1.to_binary()
    for i in range(len(gamma)-1):
        to_hash += Y[i+len(attributes_list)].to_binary()
    for gammai in gamma:
        to_hash += gammai.to_binary()

    c_hash = int(hashlib.sha512(to_hash).hexdigest(), 16)
    print("hash_test = ", c_hash)

    gamma_prod = G1.neutral_element()
    for gammai in gamma:
        gamma_prod *= gammai

    check_prod = (C**(-c_hash))*(g1**r[0])
    for i in range(len(r)-1):
        check_prod *= Y[i+len(attributes_list)]**r[i+1]
    
    # Prod(gammai) == C^c * g^r0 * Prod(Yi^ri)
    assert gamma_prod == check_prod

#@pytest.mark.skip
def test_issued_credentials(input_issuer, input_anoncred, input_cred_params):
    username = "pgab"
    G1 = input_issuer.G1
    Y = input_issuer.pk[1:]
    g = input_issuer.g
    t = input_cred_params[0]
    pub_attr_len = 2
    private_attributes = input_cred_params[1]
    C, _, _ = input_anoncred.create_issue_request(private_attributes, G1, Y, g, t, pub_attr_len)
    revealed_attr = input_cred_params[2]
    attributes = revealed_attr + private_attributes
    server_sk = Signature.deserialize(input_issuer.get_serialized_secret_key())
    sigma = input_issuer.issue(C, username, revealed_attr, server_sk)
    cred = input_anoncred.receive_issue_response(sigma, t)
    Yt = input_issuer.pkt
    gt = input_issuer.gt
    check = Yt[0]
    for i in range(len(input_issuer.pk)-1):
        check *= (Yt[i+1]**attributes[i])

    # e(sigma1, Xt*Prod(Yti^ai)) == e(sigma2, gt)
    assert cred[0].pair(check) == cred[1].pair(gt)


#@pytest.mark.skip
def test_correct_registration(input_client, input_server):
    valid_attr_str = "male, legal, private_attr"
    revealed_attr_str = "1,1"
    server_pk, server_sk = input_server.generate_ca(valid_attr_str)
    username = "pgab"
    request, prvt_state = input_client.prepare_registration(server_pk, username, revealed_attr_str)
    response = input_server.register(server_sk, request, username, revealed_attr_str)
    sigma_serial = input_client.proceed_registration_response(server_pk, response, prvt_state)
    cred = Signature.deserialize(sigma_serial).get("sigma")

    attributes = [1,1] + input_client.private_attr
    server_pk = Signature.deserialize(server_pk)
    Yt = server_pk.get("pkt")
    gt = server_pk.get("gt")
    check = Yt[0]
    for i in range(len(Yt)-1):
        check *= (Yt[i+1]**attributes[i])

    # e(sigma1, Xt*Prod(Yti^ai)) == e(sigma2, gt)
    assert cred[0].pair(check) == cred[1].pair(gt)


#@pytest.mark.skip
def test_INcorrect_registration(input_client, input_server):
    valid_attr_str = "male, legal, private_attr"
    revealed_attr_str = "1,1"
    server_pk, server_sk = input_server.generate_ca(valid_attr_str)
    username = "pgab"
    request, _ = input_client.prepare_registration(server_pk, username, revealed_attr_str)

    # incorrect commit for zkp: wrong attribute in r
    request = Signature.deserialize(request)
    r = request.get("r")
    r[1] = petrelic.bn.Bn.from_num(rd.randint(1, G1.order()))
    request = (jsonpickle.encode({"C": request.get("C"), "gamma": request.get("gamma"),
                                  "r": r})).encode('utf-8')
    response = input_server.register(server_sk, request, username, revealed_attr_str)
    assert not response

    
#@pytest.mark.skip
def test_valid_sign(input_client, input_server):
    valid_attr_str = "male, legal, private_attr"
    revealed_attr_str = "1,1"
    server_pk, server_sk = input_server.generate_ca(valid_attr_str)
    username = "pgab"
    request, prvt_state = input_client.prepare_registration(server_pk, username, revealed_attr_str)
    response = input_server.register(server_sk, request, username, revealed_attr_str)
    sigma_serial = input_client.proceed_registration_response(server_pk, response, prvt_state)
    sigma = Signature.deserialize(sigma_serial).get("sigma")
    private_attr = Signature.deserialize(sigma_serial).get("private_attr") 

    server_pk = Signature.deserialize(server_pk)
    revealed_attr = revealed_attr_str.split(',')
    message = ("i really enjoy writing those tests").encode()
    signature = Signature()
    t = signature.create_sign_request(server_pk, sigma, message, revealed_attr, private_attr)
    o = signature.sigma

    gt = server_pk.get("G2").generator()
    Yt = server_pk.get("pkt")
    Xt = Yt[0]
    Yt = Yt[1:]

    # e(o1, Xt PROD over all i Yt^ai) == e(o2, gt) 
    test = o[0].pair(Xt)
    test *= o[0].pair(gt)**t
    for i in range(len(revealed_attr)):
        test *= o[0].pair(Yt[i])**int(revealed_attr[i])
    for i in range(len(private_attr)):
        test *= o[0].pair(Yt[i + len(revealed_attr)])**int(private_attr[i])

    assert (test == (o[1].pair(gt)))


#@pytest.mark.skip
def test_verify_sign(input_client, input_server):
    valid_attr_str = "male, legal, private_attr"
    revealed_attr_str = "1,1"
    server_pk, server_sk = input_server.generate_ca(valid_attr_str)
    username = "pgab"
    request, prvt_state = input_client.prepare_registration(server_pk, username, revealed_attr_str)
    response = input_server.register(server_sk, request, username, revealed_attr_str)
    sigma_serial = input_client.proceed_registration_response(server_pk, response, prvt_state)
    sigma = Signature.deserialize(sigma_serial).get("sigma")
    private_attr = Signature.deserialize(sigma_serial).get("private_attr") 

    server_pk = Signature.deserialize(server_pk)
    revealed_attr = revealed_attr_str.split(',')
    message = ("i really enjoy debugging those tests").encode()
    signature = Signature()
    signature.create_sign_request(server_pk, sigma, message, revealed_attr, private_attr)

    assert signature.verify(server_pk, revealed_attr, message)


#@pytest.mark.skip
def test_correct_signature(input_client, input_server):
    valid_attr_str = "male, legal, private_attr"
    revealed_attr_str = "1,1"
    server_pk, server_sk = input_server.generate_ca(valid_attr_str)
    username = "pgab"
    request, prvt_state = input_client.prepare_registration(server_pk, username, revealed_attr_str)
    response = input_server.register(server_sk, request, username, revealed_attr_str)
    sigma_serial = input_client.proceed_registration_response(server_pk, response, prvt_state)
    message = ("i really enjoy not sleeping for those tests").encode()
    
    sign = input_client.sign_request(server_pk, sigma_serial, message, revealed_attr_str)
    check = input_server.check_request_signature(server_pk, message, revealed_attr_str, sign)

    assert check


#@pytest.mark.skip
def test_INcorrect_signature(input_client, input_server):
    valid_attr_str = "male, legal, private_attr"
    revealed_attr_str = "1,1"
    server_pk, server_sk = input_server.generate_ca(valid_attr_str)
    username = "pgab"
    request, prvt_state = input_client.prepare_registration(server_pk, username, revealed_attr_str)
    response = input_server.register(server_sk, request, username, revealed_attr_str)
    sigma_serial = input_client.proceed_registration_response(server_pk, response, prvt_state)
    message = ("i really enjoy not sleeping for those tests").encode()
    
    sign = input_client.sign_request(server_pk, sigma_serial, message, revealed_attr_str)

    signature = Signature.deserialize(sign)
    r = signature.r
    r[1] = petrelic.bn.Bn.from_num(rd.randint(1, G1.order()))
    signature.r = r
    wrong_sign = Signature.serialize(signature)

    check = input_server.check_request_signature(server_pk, message, revealed_attr_str, wrong_sign)

    assert not check




