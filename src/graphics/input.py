import pygame
import src.graphics.graphics_const as graphics_const


#NOTE: pygame.init() will have been called in draw.py


class Selected:
    def __init__(self, x, y, piece_value):
        #these are board coords not screen coords
        self.x = x
        self.y = y

        self.piece_value = piece_value


def get_cell_inx():
    x, y = pygame.mouse.get_pos()

    space_x = graphics_const.SCREEN_WIDTH // 8
    space_y = graphics_const.SCREEN_HEIGHT // 8

    #x and y swap because the array inx is different to cartesian coords
    cell_x = y // space_y
    cell_y = x // space_x

    return cell_x, cell_y


def player_move(board):
    global selected_piece

    x, y = get_cell_inx()

    print(x, y)

    if not 0 <= x < 8 or not 0 <= y < 8:
        return False

    if selected_piece == None:
        piece_value = board[x][y]

        #TODO: ensure piece is correct colour
        if piece_value != 0:
            piece = Selected(x, y, piece_value)
            selected_piece = piece

        #because we have not moved
        return False
    else:
        #TODO: ensure not moving onto piece of same colour

        #move piece
        board[selected_piece.x][selected_piece.y] = 0
        board[x][y] = selected_piece.piece_value

        selected_piece = None

        #because we have moved
        return True


def get_player_input(board):
    global selected_piece

    selected_piece = None

    ready_for_click = True
    while True:
        for event in pygame.event.get():
            if pygame.mouse.get_pressed()[0]:
                if ready_for_click:
                    #debounce
                    ready_for_click = False

                    has_moved = player_move(board)

                    if has_moved:
                        return board
            else:
                ready_for_click = True

            if event.type == pygame.QUIT:
                quit()