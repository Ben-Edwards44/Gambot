import src.api.api as api
import src.graphics.main as graphics
from os import system


#starting position
START_BOARD_STATE = [[10, 8, 9, 12, 11, 9, 8, 10], [7, 7, 7, 7, 7, 7, 7, 7], [0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0], [1, 1, 1, 1, 1, 1, 1, 1], [4, 2, 3, 6, 5, 3, 2, 4]]


def init():
    #call before first loop

    api.write_board_state(START_BOARD_STATE)
    graphics.init_graphics()


def main():
    #perform one loop

    board = api.load_board_state()
    player_move_board = graphics.graphics_loop(board)

    api.write_board_state(player_move_board)

    #run go engine - need to ensure the most up to date version is built
    system("chess-engine.exe")