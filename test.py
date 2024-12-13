import requests

# r = requests.get("http://127.0.0.1:8080/message")


def test_post():
    url = "http://127.0.0.1:8080/message"
    json = {
        "message1" : "ddsds",
    }

    r = requests.post(url, json = json)

    print(r.status_code)
    print(r.text)


if __name__ == "__main__":
    test_post()