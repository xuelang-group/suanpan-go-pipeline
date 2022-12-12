## 安装pip
```
pip wget https://bootstrap.pypa.io/get-pip.py
python3 get-pip.py
遇到报错No module named 'distutils.util'，
sudo apt-get install python3-distutils
python3 get-pip.py
```
### 安装fastapi
```
pip install "fastapi[all]"
```

### 运行
```
uvicorn main:app --host '0.0.0.0' --port 8080 --reload
```