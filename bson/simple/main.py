import bson


def main():
    a = bson.dumps({"a": 1, "b": [2, 3]})

    # Write the data to bson
    with open("my_file.txt", "wb") as binary_file:
        binary_file.write(a)

    # Read the data to bson
    with open('my_file.txt','rb') as f:
        b = bson.loads(f.read())
    
    print(b)


if __name__ == "__main__":
    main()