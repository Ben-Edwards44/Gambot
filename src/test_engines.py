import src.api.api as api
import src.graphics.main as graphics

from os import system
from random import randint


WHITE_PIECES = ["P", "N", "B", "R", "K", "Q"]
BLACK_PIECES = [i.lower() for i in WHITE_PIECES]

FILES = ["a", "b", "c", "d", "e", "f", "g", "h"]

FEN_FILEPATH = "data/equal_fens.txt"

SHOW_GRAPHICS = False


def parse_board(board_fen):
    ranks = board_fen.split("/")

    board = []
    for i in ranks:
        inx = 0
        rank = [0 for _ in range(8)]

        for x in i:
            if x in WHITE_PIECES:
                rank[inx] = WHITE_PIECES.index(x) + 1
                inx += 1
            elif x in BLACK_PIECES:
                rank[inx] = BLACK_PIECES.index(x) + 7
                inx += 1
            else:
                inx += int(x)

        board.append(rank)

    return board


def parse_ep(ep_target, white_to_move):
    #prev pawn double represents where the pawn is now; ep target is where the capturing pawn would go
    if ep_target == "-":
        return [-1, -1]
    
    x = int(ep_target[1])
    y = FILES.index(ep_target[0])

    if white_to_move:
        return [x + 1, y]
    else:
        return [x - 1, y]


def parse_fen(fen):
    board, move_colour, castle, ep, half_moves, full_moves = fen.split(" ")

    parsed_board = parse_board(board)

    white_to_move = move_colour == "w"

    white_k_castle = "K" in castle
    white_q_castle = "Q" in castle
    black_k_castle = "k" in castle
    black_q_castle = "q" in castle

    prev_pawn_double = parse_ep(ep, white_to_move)

    graphics.game_state.game_state_obj.load_from_dict({"board" : parsed_board,
                                                       "white_to_move" : white_to_move,
                                                       "white_king_castle" : white_k_castle,
                                                       "white_queen_castle" : white_q_castle,
                                                       "black_king_castle" : black_k_castle,
                                                       "black_queen_castle" : black_q_castle,
                                                       "prev_pawn_double" : prev_pawn_double})
    

def choose_fens(num):
    with open(FEN_FILEPATH, "r") as file:
        fens = file.read().split("\n")

    chosen = []
    for _ in range(num):
        inx = randint(0, len(fens) - 1)
        fen = fens.pop(inx)

        chosen.append(fen)

    return chosen


def run_engine(script_name):
    #run go engine

    exit_code = system(f"{script_name}.exe")
    if exit_code != 0:
        raise Exception(f"Go script ({script_name}) resulted in an error")


def check_win():
    #assumes that the graphics game state object has been updated

    api.send_data("check_win", graphics.game_state.game_state_obj)
    run_engine("chess-engine")

    terminal_state = api.load_check_win()

    return terminal_state


def play_game(white, black):
    #assumes that game state obj has been updated to the starting state of the game
    seen_pos = {}

    game_end = "not_terminal"
    while game_end == "not_terminal":
        #make move
        api.send_data("move_gen", graphics.game_state.game_state_obj)

        try:
            if graphics.game_state.game_state_obj.white_to_move:
                run_engine(white)
            else:
                run_engine(black)
        except Exception:
            game_end = "aborted"
            print(graphics.game_state.game_state_obj.board)
            quit()

        #recieve the made moved
        new_state = api.load_game_state()
        graphics.game_state.game_state_obj.load_from_dict(new_state)

        #deal with draws by repetition
        hash_board = tuple(tuple(i) for i in graphics.game_state.game_state_obj.board)            
        if hash_board in seen_pos.keys():
            seen_pos[hash_board] += 1

            if seen_pos[hash_board] >= 3:
                #draw by repetition
                game_end = "draw"
                break
        else:
            seen_pos[hash_board] = 1

        if SHOW_GRAPHICS: 
            graphics.draw_board(graphics.game_state.game_state_obj.board)
            graphics.pygame.display.update()

        game_end = check_win()

    return game_end


def engine_game(engine1, engine2, num_games):
    #play 2 different versions of the engine

    engine1 = engine1.replace("/", "\\")  #for when the scripts are in a different directory. TODO: make it work for linux
    engine2 = engine2.replace("/", "\\")  #for when the scripts are in a different directory. TODO: make it work for linux

    graphics.game_state.init_game_state(None)

    if SHOW_GRAPHICS:
        graphics.init_graphics()

    fens = choose_fens(num_games)

    win1 = 0
    win2 = 0
    draw = 0
    aborted = 0
    for i, x in enumerate(fens):
        #update game state obj
        parse_fen(x)

        #randomly assign white/black
        if randint(0, 1) == 0:
            white_player = engine1
            black_player = engine2
        else:
            white_player = engine2
            black_player = engine1

        game_end = play_game(white_player, black_player)
        
        if game_end == "draw":
            draw += 1
        elif game_end == "white_win":
            if white_player == engine1:
                win1 += 1
            else:
                win2 += 1
        elif game_end == "black_win":
            if black_player == engine1:
                win1 += 1
            else:
                win2 += 1
        else:
            aborted += 1

        print(f"Games Played: {i + 1}\n{engine1} wins: {win1}\n{engine2} wins: {win2}\nDraws: {draw}\nAborted : {aborted}\n")

    print(f"End result:\n{engine1} wins: {win1}\n{engine2} wins: {win2}\nDraws: {draw}\nAborted : {aborted}")