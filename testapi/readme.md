1. 安装pip 	
    wget https://bootstrap.pypa.io/get-pip.py
    python3 get-pip.py
    遇到报错No module named 'distutils.util'，按下面方式按爪个包
    sudo apt-get install python3-distutils
    sudo apt-get install python3-apt
2. 安装fastapi: (https://fastapi.tiangolo.com/zh/tutorial/), pip install "fastapi[all]"
3. 修改默认端口启动fastapi: vicorn main:app --host '192.168.30.17' --port 8000 --reload