import draw
import graphics_const

import pygame
from time import time


class Player:
    def __init__(self, board, colour):
        self.board = board
        self.colour = colour


class EnginePlayer(Player):
    def __init__(self, board, colour):
        super().__init__(board, colour)

    def get_move(self):
        #get the player's move - assumes position has been updated
        if graphics_const.USE_CLOCK_TIME:
            args = {"wtime" : self.board.white_clock.get_ms(), "btime" : self.board.black_clock.get_ms()}
        else:
            args = {"movetime" : graphics_const.ENGINE_MOVE_TIME}

        self.board.engine.send_args("go", args)

        move = self.board.engine.try_get_move()
        while move == None:
            draw.draw_clocks(self.board)
            move = self.board.engine.try_get_move()

        return move


class HumanPlayer(Player):
    def __init__(self, board, colour):
        super().__init__(board, colour)

        self.dragging_piece = False
        self.dragging_piece_x = None
        self.dragging_piece_y = None

    def reset_dragging(self):
        self.dragging_piece = False
        self.dragging_piece_x = None
        self.dragging_piece_y = None

    def kill_process(self):
        self.board.kill_process()

    def check_quit(self):
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                self.board.kill_process()
                quit()

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

    def convert_move(self, start_x, start_y, end_x, end_y):
        #convert the dragging coords to a UCI move
        prom_value = self.check_for_promotion(start_x, start_y, end_x)  #promotions are special case
        move = self.board.move_to_str(start_x, start_y, end_x, end_y, prom_value)

        return move
    
    def check_legal(self, start_x, start_y, end_x, end_y):
        #make sure the player's move is legal
        if not graphics_const.LEGAL_FILTER:
            return True

        target = (end_x, end_y)
        legal_moves = self.board.get_legal_moves(start_x, start_y)

        return target in legal_moves


    def check_for_promotion(self, start_x, start_y, end_rank):
        piece_val = self.board.board_list[start_x][start_y]

        if (piece_val == 1 or piece_val == 7) and (end_rank == 0 or end_rank == 7):
            #promotion
            draw.draw_promotion(end_rank, start_y, piece_val)
            
            promotion_inx = self.get_promoted_piece(end_rank, start_y)
            promotion_val = graphics_const.PROMOTION_ORDER[promotion_inx]
        else:
            promotion_val = ""

        return promotion_val


    def get_promoted_piece(self, pawn_x, pawn_y):
        #get the player's choice of promotion from the menu
        chosen_inx = None
        while chosen_inx is None:
            if pygame.mouse.get_pressed()[0]:
                x, y = self.get_mouse_coords()
                vert_offset = abs(pawn_x - x)

                if y == pawn_y and vert_offset < 4:  #if clicked on menu
                    chosen_inx = vert_offset

            draw.draw_clocks(self.board)
            self.check_quit()

        return chosen_inx


    def get_move(self):
        #get the player's move - assumes position has been updated
        made_move = False

        while not made_move:
            clicked = pygame.mouse.get_pressed()[0]

            draw.draw_clocks(self.board)  #we only need to draw the clocks here because the board may be unchanged

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

            self.check_quit()

        return move