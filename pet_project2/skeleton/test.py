from credential import PSSignature

machin = PSSignature()
sk, pk = machin.generate_key()
print("secret is ", sk[0], type(sk[0]))
print("---------------------------------")
signature = machin.sign(sk, [b"Fake news", b"apchii"])
print('signature is ', signature)

print("---------------------------------")
print(machin.verify(pk, signature, [b"Fake news", b"apchii"]))