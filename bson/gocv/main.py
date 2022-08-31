import bson
import numpy as np 
import cv2

from typing import ByteString

def main():
    # Read the data to bson
    with open('matrix','rb') as f:
        b = bson.loads(f.read())
        print(b.keys())
        nparr = np.asarray(bytearray(b["data"]), dtype=np.uint8).reshape((2521, 3361, 3))
        cv2.imwrite("./mugFromPython.jpg", nparr)

if __name__ == '__main__':
    main()
