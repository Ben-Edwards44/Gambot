import json


FILE_PATH = "src/api/interface.json"


def str_to_list(string):
    string = string[2:-2].split("], [")
    
    list = []
    for i in string:
        nums = i.split(", ")
        nums = [int(x) for x in nums]

        list.append(nums)

    return list


def load_board_state():
    with open(FILE_PATH, "r") as file:
        data = file.read()

    #board_dict in form {"board" : "[[...], [...], ...]"}
    board_dict = json.loads(data)
    board_str = board_dict["board"]

    board = str_to_list(board_str)

    return board


def write_board_state(board):
    #board in form [[x, y, z, ...], [...], ...] where 0 = empty, 1 = white pawn etc.

    board_dict = {"board" : str(board)}
    json_str = json.dumps(board_dict)

    with open(FILE_PATH, "w") as file:
        file.write(json_str)