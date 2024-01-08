import draw
import piece
from os import listdir


def init_graphics():
    global pieces

    draw.init_draw()

    images = get_images()
    pieces = build_pieces(images)


def get_images():
    images = []
    for i in ["white", "black"]:
        path = f"images\\{i}"
        img_names = listdir(path)

        #images in form {name : path} e.g. {"queen" : "Images/White/queen.png"}
        images.append({x[:-4] : f"{path}\\{x}" for x in img_names})

    return images


def build_pieces(images):
    #return list of Piece objects with starting position and image path

    white_row = ["rook", "knight", "bishop", "queen", "king", "bishop", "knight", "rook"]
    black_row = white_row[::-1]

    x_pos = 7
    pieces = []

    #x_pos refers to 1st index of 2d array, y_pos refers to 2nd (opposite to cartesian coords)
    for x_pos in [0, 1, 6, 7]:
        is_white = x_pos > 4
        is_pawn = 0 < x_pos < 7

        row = white_row if is_white else black_row
        img_dict = images[0] if is_white else images[1]

        for y_pos in range(8):
            if is_pawn:
                name = "pawn"
            else:
                name = row[y_pos]

            img_path = img_dict[name]
            new_piece = piece.Piece(name, img_path, x_pos, y_pos)

            pieces.append(new_piece)

    return pieces


def draw_board():
    #draw background and pieces
    draw.draw_board(pieces)


if __name__ == "__main__":
    #for testing

    init_graphics()

    draw_board()

    while True:
        pass