import src.api.api as api
import src.graphics.main as graphics
import src.graphics.game_state as game_state

from os import system


def temp(a):
    a = [int(i) for i in a.split(" ")]
    n = [[0 for _ in range(8)] for _ in range(8)]

    for i in range(8):
        for x in range(8):
            n[i][x] = a[i * 8 + x]

    return n


#starting position
START_BOARD_STATE = [[10, 8, 9, 12, 11, 9, 8, 10], [7, 7, 7, 7, 7, 7, 7, 7], [0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0], [1, 1, 1, 1, 1, 1, 1, 1], [4, 2, 3, 6, 5, 3, 2, 4]]
PLAYER_WHITE = True


def init():
    #call before first loop

    game_state.init_game_state(START_BOARD_STATE)
    graphics.init_graphics()

    api.send_data("move_gen", game_state.game_state_obj)


def run_engine():
    #run go engine

    exit_code = system("chess-engine.exe")
    if exit_code != 0:
        raise Exception("Go script resulted in an error")


def perft(depth):
    #do performance test

    game_state.init_game_state(START_BOARD_STATE)
    api.send_data("perft", game_state.game_state_obj, perft_depth=depth)

    run_engine()


def main():
    #perform one loop

    state_dict = api.load_game_state()
    graphics.game_state.game_state_obj.load_from_dict(state_dict)
    graphics.game_state.game_state_obj.white_to_move = PLAYER_WHITE  #because it is player's turn

    player_move_board = graphics.graphics_loop(graphics.game_state.game_state_obj.board)

    #ensure the game state is updated
    graphics.game_state.game_state_obj.board = player_move_board
    graphics.game_state.game_state_obj.white_to_move = not PLAYER_WHITE  #because it is no longer player turn

    api.send_data("move_gen", graphics.game_state.game_state_obj)

    run_engine()