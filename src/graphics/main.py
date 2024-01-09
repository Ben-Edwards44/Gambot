import src.graphics.draw as draw
import src.graphics.piece as piece
import src.graphics.graphics_const as graphics_const

from os import listdir


def init_graphics():
    global images

    draw.init_draw()

    images = get_images()


def get_images():
    images = []
    for i in ["white", "black"]:
        path = f"src/graphics/images/{i}"
        img_names = listdir(path)

        #images in form {name : path} e.g. {"queen" : "Images/White/queen.png"}
        images.append({x[:-4] : f"{path}\\{x}" for x in img_names})

    return images


def build_pieces(board, images):
    #return list of Piece objects with starting position and image path

    piece_list = []
    for x in range(8):
        for y in range(8):
            value = board[x][y]

            if value != 0:
                #find colour
                if value > 6:
                    is_white = False
                    value -= 6
                else:
                    is_white = True

                name = graphics_const.PIECE_VALUES[value]

                #get image
                img_dict = images[0] if is_white else images[1]
                img_path = img_dict[name]

                #add piece to list
                new_piece = piece.Piece(name, img_path, x, y)
                piece_list.append(new_piece)
            
    return piece_list


def draw_board(board):
    #convert board to list of pieces
    piece_list = build_pieces(board, images)

    #actually draw the board and background
    draw.draw_board(piece_list)