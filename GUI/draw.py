import piece
import graphics_const

import pygame


class GraphicsBoard:
    def __init__(self, board, pieces, dragging_piece=None):
        self.board = board
        self.pieces = pieces
        self.dragging_piece = dragging_piece

        self.font = pygame.font.Font(graphics_const.FONT_NAME, graphics_const.FONT_SIZE)

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

    def draw_text(self, text, x, y):
        text_surface = self.font.render(text, graphics_const.FONT_ANTIALIAS, graphics_const.FONT_COLOUR)
        text_rect = text_surface.get_rect()

        text_rect.center = (x, y)

        window.blit(text_surface, text_rect)

    def draw_border(self):
        pygame.draw.rect(window, graphics_const.BORDER_COLOUR, (graphics_const.BORDER_TL[0], graphics_const.BORDER_TL[1], graphics_const.BORDER_X, graphics_const.BORDER_Y), border_radius=graphics_const.BORDER_CORNER_RADIUS)

        #draw the file labels a, b, c, d etc.
        file_y = graphics_const.BOARD_TL[1] + graphics_const.BOARD_Y + graphics_const.BORDER_WIDTH // 2
        for i, file in enumerate(graphics_const.FILES):
            file_x = graphics_const.BOARD_TL[0] + i * graphics_const.STEP_X + graphics_const.STEP_X // 2

            self.draw_text(file, file_x, file_y)

        #draw the rank labels 1, 2, 3, 4 etc.
        rank_x = graphics_const.BOARD_TL[0] - graphics_const.BORDER_WIDTH // 2
        for i in range(8):
            rank = f"{8 - i}"
            rank_y = graphics_const.BOARD_TL[1] + i * graphics_const.STEP_Y + graphics_const.STEP_Y // 2

            self.draw_text(rank, rank_x, rank_y)

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