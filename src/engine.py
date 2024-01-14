import src.api.api as api
import src.graphics.main as graphics
import src.graphics.game_state as game_state

from os import system


#starting position
START_BOARD_STATE = [[10, 8, 9, 12, 11, 9, 8, 10], [7, 7, 7, 7, 7, 7, 7, 7], [0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0], [1, 1, 1, 1, 1, 1, 1, 1], [4, 2, 3, 6, 5, 3, 2, 4]]


def init():
    #call before first loop

    game_state.init_game_state(START_BOARD_STATE)
    graphics.init_graphics()

    api.send_data("move_gen", game_state.game_state_obj)


def main():
    #perform one loop

    state_dict = api.load_game_state()
    graphics.game_state.game_state_obj.load_from_dict(state_dict)

    player_move_board = graphics.graphics_loop(graphics.game_state.game_state_obj.board)

    #ensure the game state is updated
    graphics.game_state.game_state_obj.board = player_move_board

    api.send_data("move_gen", graphics.game_state.game_state_obj)

    #run go engine - need to ensure the most up to date version is built
    exit_code = system("chess-engine.exe")
    if exit_code != 0:
        raise Exception("Go script resulted in an error")