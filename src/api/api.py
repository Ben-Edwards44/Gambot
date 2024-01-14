import json


FILE_PATH = "src/api/interface.json"


def split_elements(string):
    elements = []

    net_brackets = 0
    current_element = ""

    for i in string:
        if i == "," or i == " ":
            if net_brackets == 0:
                if current_element != "":
                    elements.append(current_element)
                    current_element = ""
            else:
                current_element += i
        elif i == "[":
            net_brackets += 1
            current_element += i
        elif i == "]":
            net_brackets -= 1
            current_element += i
        else:
            current_element += i

    if current_element != "":
        elements.append(current_element)

    return elements


def str_to_list(string):
    #recursively parse nested lists
    
    open_inx = None
    closed_inx = None

    l = len(string)
    for i in range(l):
        if string[i] == "[" and open_inx == None:
            open_inx = i
        if string[l - i - 1] == "]" and closed_inx == None:
            closed_inx = l - i - 1

    string = string[open_inx + 1 : closed_inx]
    elements = split_elements(string)

    list = []
    for i in elements:
        if "[" in i:
            parsed = str_to_list(i)
        else:
            parsed = int(i)

        list.append(parsed)

    return list


def read_json():
    with open(FILE_PATH, "r") as file:
        data = file.read()

    data_dict = json.loads(data)

    return data_dict


def parse_json(json):
    parsed = {}
    for k, v in json.items():
        if v[0] == "[":
            #works for board and prev_pawn_double
            parsed[k] = str_to_list(v)
        elif v == "true" or v == "false":
            parsed[k] = v == "true"
        else:
            parsed[k] = v

    return parsed


def load_game_state():
    data_dict = read_json()
    parsed = parse_json(data_dict)

    return parsed


def load_legal_moves():
    move_dict = read_json()

    moves = move_dict["moves"]
    moves = str_to_list(moves)

    return moves


def concat_dict(og_dict, add_dict):
    for k, v in add_dict.items():
        #go side of api needs everything as string
        t = type(v)

        if t == bool:
            v = "true" if v else "false"
        elif t != str:
            v = str(v)

        og_dict[k] = v


def send_data(engine_task, game_state_obj, **kwargs):
    #game_state_obj has attrs like board, castle data etc.

    send_dict = {"task" : engine_task}

    concat_dict(send_dict, game_state_obj.get_dict())
    concat_dict(send_dict, kwargs)

    json_str = json.dumps(send_dict)

    with open(FILE_PATH, "w") as file:
        file.write(json_str)