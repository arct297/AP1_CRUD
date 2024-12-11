import requests

r = requests.get("http://127.0.0.1:8080/message")

print(r.status_code)
print(r.text)