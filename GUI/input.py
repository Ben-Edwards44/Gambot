import draw
import graphics_const

import pygame


def get_cell_inx():
    x, y = pygame.mouse.get_pos()

    space_x = graphics_const.SCREEN_WIDTH // 8
    space_y = graphics_const.SCREEN_HEIGHT // 8

    #x and y swap because the array inx is different to cartesian coords
    cell_x = y // space_y
    cell_y = x // space_x

    return cell_x, cell_y


def convert_move(selected_x, selected_y):
    #get a move string of the player's move

    end_x, end_y = get_cell_inx()

    start_file = graphics_const.FILES[selected_y]
    end_file = graphics_const.FILES[end_y]

    move = f"{start_file}{8 - selected_x}{end_file}{8 - end_x}"

    return move


def drag_piece(board, selected_x, selected_y):
    if selected_x == -1 and selected_y == -1:
        #player has just selected a piece
        x, y = get_cell_inx()

        #TODO: also check if it is correct colour
        if board[x][y] != 0:
            selected_x = x
            selected_y = y
    else:
        #player is moving selected piece
        draw.draw_dragging_piece(board, selected_x, selected_y)

    return selected_x, selected_y


def get_move(board):
    #get the player's move

    selected_x = -1
    selected_y = -1

    while True:        
        if pygame.mouse.get_pressed()[0]:
            selected_x, selected_y = drag_piece(board, selected_x, selected_y)
        else:
            if selected_x != -1 and selected_y != -1:
                move = convert_move(selected_x, selected_y)
                #TODO: check if move in legal moves (just get a list of all moves from the chess engine)

                return move

            selected_x = -1
            selected_y = -1

        #need to pump to ensure clicks are properly handeled
        pygame.event.pump()