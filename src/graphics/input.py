import pygame
import src.graphics.graphics_const as graphics_const


#NOTE: pygame.init() will have been called in graphics/main.py


selected_piece = None


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


def select(board):
    global selected_piece

    if selected_piece != None:
        return

    x, y = get_cell_inx()
    piece_value = board[x][y]

    #TODO: ensure piece is correct colour
    if piece_value != 0:
        piece = Selected(x, y, piece_value)
        selected_piece = piece


def move_selected(board):
    global selected_piece

    if selected_piece == None:
        return
    
    x, y = get_cell_inx()
    
    board[selected_piece.x][selected_piece.y] = 0
    board[x][y] = selected_piece.piece_value
    selected_piece = None


def get_player_input(board):
    global selected_piece

    #need to pump to ensure clicks are properly handeled
    pygame.event.pump()

    if pygame.mouse.get_pressed()[0]:
        select(board)
    else:
        move_selected(board)

    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            quit()

    return board