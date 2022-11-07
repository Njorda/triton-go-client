# Run the go code


From the root folder, to run the go code: 
```bash
docker run -it -v $(pwd):/go/src gocv/opencv /bin/bash
```
and then: 

```bash
cd /go/src/bson/gocv
go run main.go
```



From the root folder, to run the python code:

```python
docker run -it -v $(pwd):/home opencvcourses/opencv-docker /bin/bash
```

install deps
```
pip install bson
```

open the bson in python
```
cd /bson/gocv
python main.py
```