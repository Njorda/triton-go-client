
import bson
import io
import torch

# triton_python_backend_utils is available in every Triton Python model. You
# need to use this module to create inference requests and responses. It also
# contains some utility functions for extracting information from model_config
# and converting Triton input/output types to numpy types.
import torchvision.transforms as transforms
import numpy as np
from PIL import Image
from matplotlib import pyplot as plt



from typing import List

def main():
    normalize = transforms.Normalize(mean=[0.485, 0.456, 0.406],
                                                std=[0.229, 0.224, 0.225])

    loader = transforms.Compose([
        transforms.Resize([224, 224]),
        transforms.CenterCrop(224),
        transforms.ToTensor(), normalize
    ])

    def image_loader(image_name):
        image = loader(image_name)
        image = image.unsqueeze(0)
        return image

    with open('matrix','rb') as f:
        b = bson.loads(f.read())

    img = np.asarray(bytearray(b["data"]), dtype=np.uint8).reshape((2521, 3361, 3)) 
    im_pil = Image.fromarray(img)
    img_out = image_loader(im_pil)
    img_out = np.array(img_out)
    print(img_out)
    

if __name__ == '__main__':
    main()
