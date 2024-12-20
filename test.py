import requests


def _send_request(endpoint, json, rtype="post"):
    url = f"http://127.0.0.1:8080{endpoint}"

    print(f"request_type={rtype} json={json} {url=}")
    r = requests.request(rtype, url, json = json)

    received_status_code = r.status_code
    try:
        received_json = r.json()
    except:
        received_json = {}

    return received_status_code, received_json


def test_task_1():    
    endpoint = "/message"

    # Test 1
    # Correct message
    field_name = "message"
    field_content = "whtrmtkrnkrtmkr"
    json = {
        field_name : field_content
    }
    received_status_code, received_json = _send_request(endpoint, json)
    print(f"Test 1, correct message.\nrequest data: {field_name} : {field_content};\nReceived status code: {received_status_code}; received json response: {received_json}\n")

    # Test 2
    # Incorrect message, invalid type
    field_name = "message"
    field_content = 1
    json = {
        field_name : field_content
    }
    received_status_code, received_json = _send_request(endpoint, json)
    print(f"Test 2, incorrect message.\nrequest data: {json=};\nReceived status code: {received_status_code}; received json response: {received_json}\n")

    # Test 3
    # No expected message field, another field
    field_name = "message1"
    field_content = "Waiting error..."
    json = {
        field_name : field_content
    }
    received_status_code, received_json = _send_request(endpoint, json)
    print(f"Test 3, incorrect message field.\nrequest data: {field_name} : {field_content};\nReceived status code: {received_status_code}; received json response: {received_json}\n")
    

def test_task_2():
    endpoint = "/patients"
    json = {
        "name" : "Andr",
        "age" : 21,
        "gender" : "Male",
        "contact" : "+77000000000",
        "address" : "Kazakhstan state,..."
    }
    received_status_code, received_json = _send_request(endpoint, json)
    print(received_status_code, received_json)

    # id = received_json.get("ID")

    id = 5

    # # GET BY ID
    # endpoint = f"/patients/{id}"
    # received_status_code, received_json = _send_request(endpoint, {}, "get")
    # print(received_status_code, received_json)

    # # UPDATE BY ID
    # endpoint = f"/patients/{id}"
    # json = {
    #     "age" : 20
    # }
    # # json = {}
    # received_status_code, received_json = _send_request(endpoint, json, "put")
    # print(received_status_code, received_json)

    # # GET BY ID
    # endpoint = f"/patients/{id}"
    # received_status_code, received_json = _send_request(endpoint, {}, "get")
    # print(received_status_code, received_json)

    # GET BY ID
    # endpoint = f"/patients/{id}"
    # received_status_code, received_json = _send_request(endpoint, {}, "get")
    # print(received_status_code, received_json)

    # # DELETE BY ID
    # endpoint = f"/patients/{id}"
    # received_status_code, received_json = _send_request(endpoint, {}, "delete")
    # print(received_status_code, received_json)

    # # GET BY ID
    # endpoint = f"/patients/{id}"
    # received_status_code, received_json = _send_request(endpoint, {}, "get")
    # print(received_status_code, received_json)


if __name__ == "__main__":
    # test_task_1()
    test_task_2()
