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


def read_json():
    with open(FILE_PATH, "r") as file:
        data = file.read()

    data_dict = json.loads(data)

    return data_dict


def load_board_state():
    board_dict = read_json()

    #board_dict in form {"board" : "[[...], [...], ...]", ...}
    board_str = board_dict["board"]

    board = str_to_list(board_str)

    return board


def load_legal_moves():
    move_dict = read_json()

    moves = move_dict["moves"]
    moves = str_to_list(moves)

    return moves


def send_data(engine_task, board, **kwargs):
    #board in form [[x, y, z, ...], [...], ...] where 0 = empty, 1 = white pawn etc.

    send_dict = {"task" : engine_task, "board" : str(board)}

    for name, val in kwargs.items():
        #go side of api needs everything as string
        if type(val) != str:
            val = str(val)

        send_dict[name] = val

    json_str = json.dumps(send_dict)

    with open(FILE_PATH, "w") as file:
        file.write(json_str)