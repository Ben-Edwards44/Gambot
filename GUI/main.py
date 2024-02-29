import draw
import input
import graphics_const
import engine_interface

import pygame


START_BOARD = [[10, 8, 9, 12, 11, 9, 8, 10], [7, 7, 7, 7, 7, 7, 7, 7], [0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0], [0, 0, 0, 0, 0, 0, 0, 0], [1, 1, 1, 1, 1, 1, 1, 1], [4, 2, 3, 6, 5, 3, 2, 4]]


def make_move(move, board):
    #NOTE: need to do castling / ep / promotions etc.
    start_x = 8 - int(move[1])
    start_y = graphics_const.FILES.index(move[0])
    end_x = 8 - int(move[3])
    end_y = graphics_const.FILES.index(move[2])

    piece_val = board[start_x][start_y]

    board[start_x][start_y] = 0
    board[end_x][end_y] = piece_val


def player_move(board):
    move = input.get_move(board)

    return move
        

def engine_move(engine, move_list):
    engine.set_pos(move_list)

    return engine.get_move()
        

def main():
    draw.init()

    move_list = []
    white_to_move = True
    engine = engine_interface.Engine()

    board = START_BOARD

    draw.draw_board(board)

    #TODO: do isready, ucinewgame etc.

    while True:
        if white_to_move == graphics_const.PLAYER_WHITE:
            move = player_move(board)
        else:
            move = engine_move(engine, move_list)

        print(move)

        make_move(move, board)

        move_list.append(move)
        white_to_move = not white_to_move

        draw.draw_board(board)

        for e in pygame.event.get():
            if e.type == pygame.QUIT:
                del engine

                quit()


if __name__ == "__main__":
    main()