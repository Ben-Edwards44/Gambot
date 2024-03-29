import piece
import graphics_const

import pygame


class GraphicsBoard:
    def __init__(self, board, pieces, dragging_piece=None):
        self.board = board
        self.pieces = pieces
        self.dragging_piece = dragging_piece

    def draw_square(self, x, y, colour):
        draw_x = x * graphics_const.STEP_X + graphics_const.BOARD_TL[0]
        draw_y = y * graphics_const.STEP_Y + graphics_const.BOARD_TL[1]

        pygame.draw.rect(window, colour, (draw_x, draw_y, graphics_const.STEP_X, graphics_const.STEP_Y))

    def draw_background_squares(self):
        #draw the background squares
        for i in range(8):
            for j in range(8):
                if (i + j) % 2 == 0:
                    colour = graphics_const.LIGHT_SQ_COLOUR
                else:
                    colour = graphics_const.DARK_SQ_COLOUR

                self.draw_square(i, j, colour)

    def draw_border(self):
        pygame.draw.rect(window, graphics_const.BORDER_COLOUR, (graphics_const.BORDER_TL[0], graphics_const.BORDER_TL[1], graphics_const.BORDER_X, graphics_const.BORDER_Y))

    def draw_legal_moves(self, x, y):
        #colour the squares that the player could move to
        legal_moves = self.board.get_legal_moves(x, y)

        for end_x, end_y in legal_moves:
            draw_x = end_y
            draw_y = end_x

            self.draw_square(draw_x, draw_y, graphics_const.LEGAL_MOVE_COLOUR)

    def draw(self):
        self.draw_border()
        self.draw_background_squares()

        if self.dragging_piece is not None:
            self.draw_legal_moves(self.dragging_piece.x, self.dragging_piece.y)

        for i in self.pieces:
            i.draw()

        if self.dragging_piece is not None:
            self.dragging_piece.draw()  #draw the dragging piece last so that it appears on top of any other pieces


def init():
    global window

    pygame.init()
    pygame.display.set_caption("Chess Engine")

    window = pygame.display.set_mode((graphics_const.SCREEN_WIDTH, graphics_const.SCREEN_HEIGHT))


def draw_board(board, dragging_x=None, dragging_y=None):
    pieces_list, dragging_piece = piece.get_pieces(window, board.board_list, dragging_x, dragging_y)

    board = GraphicsBoard(board, pieces_list, dragging_piece)
    board.draw()

    pygame.display.update()