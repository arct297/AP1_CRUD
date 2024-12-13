import requests


def _send_request(json):
    url = "http://127.0.0.1:8080/message"

    r = requests.post(url, json = json)

    received_status_code = r.status_code
    try:
        received_json = r.json()
    except:
        received_json = {}

    return received_status_code, received_json


def test_task_1():    
    # Test 1
    # Correct message
    field_name = "message"
    field_content = "Hi! Waiting success..."
    json = {
        field_name : field_content
    }
    received_status_code, received_json = _send_request(json)
    print(f"Test 1, correct message.\nrequest data: {field_name} : {field_content};\nReceived status code: {received_status_code}; received json response: {received_json}\n")

    # Test 2
    # Incorrect message, invalid type
    field_name = "message"
    field_content = 1
    json = {
        field_name : field_content
    }
    received_status_code, received_json = _send_request(json)
    print(f"Test 2, incorrect message.\nrequest data: {json=};\nReceived status code: {received_status_code}; received json response: {received_json}\n")

    # Test 3
    # No expected message field, another field
    field_name = "message1"
    field_content = "Waiting error..."
    json = {
        field_name : field_content
    }
    received_status_code, received_json = _send_request(json)
    print(f"Test 3, incorrect message field.\nrequest data: {field_name} : {field_content};\nReceived status code: {received_status_code}; received json response: {received_json}\n")
    


if __name__ == "__main__":
    test_task_1()