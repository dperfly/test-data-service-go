import requests
import jsonpath

IP = "127.0.0.1"
Port = "8080"


class GetDataExample:
    def txt_data_example(self):
        file_name = "txt_demo.txt"
        res = requests.get(f"http://{IP}:{Port}/{file_name}?num=3")
        json_list = jsonpath.jsonpath(res.json(), "$..result[*]")
        print(json_list)
        # ['18666660004', '18666660005', '18666660006']

    def csv_data_example(self):
        file_name = "csv_demo.csv"
        res = requests.get(f"http://{IP}:{Port}/{file_name}?num=3")
        json_list = jsonpath.jsonpath(res.json(), "$..result[*]")
        print(json_list)
        # [['user1', '123456'], ['user2', '123456'], ['user3', '123456']]

    def json_data_example(self):
        file_name = "json_demo.json"
        res = requests.get(f"http://{IP}:{Port}/{file_name}?num=3")
        result_data = res.json()["result"]
        print(result_data)
        # {'age': 18, 'name': 'zhang san'}

