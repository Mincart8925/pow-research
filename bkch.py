import hashlib
from binascii import unhexlify, hexlify

class Blocks:
    version = "01000000"
    prevhash = "00000000000008a3a41b85b8b29ad444def299fee21793cd8b9e567eab02cd81"
    hashmerkle = "2b12fcf1b09288fcaff797d71e950e71ae42b91e8bdb2304758dfcffc2b620e3"
    timestamp = "1305998791"
    diff = "f2b9441a"
    nonce = "00000000"
    maked_diff = ""
    firsthex = ""
    block_hash = ""

    def __init__(self):
        print("Block Created!")
        
    def get_block_hash(self):
        tmphash = ( self.version + self.prevhash + self.hashmerkle + self.timestamp + self.diff + self.nonce )
        tmphash2 = unhexlify(tmphash)
        
        tmphash3 = hashlib.sha256(hashlib.sha256(tmphash2).digest()).digest()
        self.block_hash = hexlify(tmphash3).decode("utf-8")

        return self.block_hash

    def mine_block(self):
        print("Mining blocks...")
        print("Difficulty is: ", int(self.diff, 8))
        self.make_diff_str()
        #while not 

    def make_diff_str(self):
        for _ in range(0, int(self.diff)):
            self.maked_diff += "0"

bk = Blocks()
data = bk.get_block_hash()
print("Block hash: ", data)
bk.mine_block()
        
        
