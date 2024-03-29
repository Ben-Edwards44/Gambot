import draw
import graphics_const

import pygame
from time import time


class Player:
    def __init__(self, board, colour, start_time, time_allowed):
        self.board = board
        self.colour = colour
        self.start_time = start_time
        self.time_allowed = time_allowed

    get_time_left = lambda self: time() - self.start_time


class EnginePlayer(Player):
    def __init__(self, board, colour, start_time, time_allowed):
        super().__init__(board, colour, start_time, time_allowed)

    def get_move(self):
        #get the player's move - assumes position has been updated
        move = self.board.engine.get_move(movetime=graphics_const.ENGINE_MOVE_TIME)

        return move


class HumanPlayer(Player):
    def __init__(self, board, colour, start_time, time_allowed):
        super().__init__(board, colour, start_time, time_allowed)

        self.dragging_piece = False
        self.dragging_piece_x = None
        self.dragging_piece_y = None

    def reset_dragging(self):
        self.dragging_piece = False
        self.dragging_piece_x = None
        self.dragging_piece_y = None

    def kill_process(self):
        self.board.kill_process()

    def get_mouse_coords(self):
        #get the board cell coordinates of the mouse position
        x, y = pygame.mouse.get_pos()

        #adjust for the fact the the board is offset
        x -= graphics_const.BOARD_TL[0]
        y -= graphics_const.BOARD_TL[1]

        #x and y swap because the array inx is different to cartesian coords
        cell_x = y // graphics_const.STEP_Y
        cell_y = x // graphics_const.STEP_X

        return cell_x, cell_y

    def start_dragging_piece(self):
        x, y = self.get_mouse_coords()

        piece_val = self.board.board_list[x][y]

        if piece_val != 0:
            self.dragging_piece = True
            self.dragging_piece_x = x
            self.dragging_piece_y = y

    def draw_dragging_piece(self):
        #draw the piece as we are dragging it
        if self.dragging_piece:
            draw.draw_board(self.board, self.dragging_piece_x, self.dragging_piece_y)

    def check_for_promotion(self, end_x, end_y):
        piece_val = self.board.board_list[end_x][end_y]

        if (piece_val == 1 or piece_val == 7) and (end_x == 0 or end_x == 7):
            #promotion
            promotion_val = input("Enter promotion piece: ")
        else:
            promotion_val = ""

        return promotion_val

    def convert_move(self, start_x, start_y, end_x, end_y):
        #convert the dragging coords to a UCI move
        prom_value = self.check_for_promotion(end_x, end_y)  #promotions are special case
        move = self.board.move_to_str(start_x, start_y, end_x, end_y, prom_value)

        return move
    
    def check_legal(self, start_x, start_y, end_x, end_y):
        #make sure the player's move is legal
        target = (end_x, end_y)
        legal_moves = self.board.get_legal_moves(start_x, start_y)

        return target in legal_moves


    def get_move(self):
        #get the player's move - assumes position has been updated
        made_move = False

        while not made_move:
            clicked = pygame.mouse.get_pressed()[0]

            if clicked:
                if not self.dragging_piece:
                    self.start_dragging_piece()
                else:
                    self.draw_dragging_piece()
            else:
                if self.dragging_piece:
                    #the player has just released a piece after dragging
                    end_x, end_y = self.get_mouse_coords()
                    move = self.convert_move(self.dragging_piece_x, self.dragging_piece_y, end_x, end_y)

                    if self.check_legal(self.dragging_piece_x, self.dragging_piece_y, end_x, end_y):
                        made_move = True
                    else:
                        draw.draw_board(self.board)  #the board must be redrawn otherwise the dragging piece will just float

                    self.reset_dragging()

            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    self.board.kill_process()
                    quit()

        return move