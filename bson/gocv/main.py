import bson
import numpy as np 
import cv2

from typing import ByteString


def decode(raw: ByteString):
    b = np.fromstring(raw, np.uint8).reshape(3361,2521, 3)
    print(b)
    # it is empty and that is the problem ...
    img = cv2.imdecode(b, flags=cv2.IMREAD_COLOR)
    print(img)
    cv2.imwrite("out.jpg", img)

def main():
    # Read the data to bson
    print("START")
    print("START")
    print("START")
    print("START")
    print("START")
    with open('matrix.txt','rb') as f:
        b = bson.loads(f.read())
    print("HERE WE GO GO GO GO ")
    print("HERE WE GO GO GO GO ")
    print("HERE WE GO GO GO GO ")
    print(b.keys())
    decode(b["data"])

if __name__ == '__main__':
    main()
