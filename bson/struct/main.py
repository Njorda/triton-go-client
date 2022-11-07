import bson
import numpy as np 
import cv2

def main():
    # Read the data to bson
    with open('matrixStruct','rb') as f:
        b = bson.loads(f.read())
        print(b.keys())
        nparr = np.asarray(bytearray(b["image"]), dtype=np.uint8).reshape((2521, 3361, 3))
        cv2.imwrite(b["id"], nparr)
        print("TOW TOW")

if __name__ == '__main__':
    main()