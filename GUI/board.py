import graphics_const


class Board:
    def __init__(self, fen, engine_process):
        self.start_fen = fen
        self.engine = engine_process

        self.move_list = []

        self.board_list = self.fen_to_board(fen)

    def kill_process(self):
        self.engine.kill_process()

    def str_to_square(self, square):
        x = 8 - int(square[1])
        y = graphics_const.FILES.index(square[0])

        return x, y

    def square_to_str(self, x, y):
        rank = 8 - x
        file = graphics_const.FILES[y]

        return f"{file}{rank}"
    
    def move_to_str(self, start_x, start_y, end_x, end_y, promotion_val):
        start = self.square_to_str(start_x, start_y)
        end = self.square_to_str(end_x, end_y)

        return f"{start}{end}{promotion_val}"

    def str_to_move(self, move):
        start_x, start_y = self.str_to_square(move[:2])
        end_x, end_y = self.str_to_square(move[2:])

        return start_x, start_y, end_x, end_y

    def fen_to_board(self, fen):
        #return a 2d self.board_list array from a fen string
        b_fen = fen.split(" ")[0]
        ranks = b_fen.split("/")

        board_list = []
        for i in ranks:
            inx = 0
            rank = [0 for _ in range(8)]

            for x in i:
                if x in graphics_const.WHITE_PIECES:
                    rank[inx] = graphics_const.WHITE_PIECES.index(x) + 1
                    inx += 1
                elif x in graphics_const.BLACK_PIECES:
                    rank[inx] = graphics_const.BLACK_PIECES.index(x) + 7
                    inx += 1
                else:
                    inx += int(x)

            board_list.append(rank)

        return board_list
    
    def make_move(self, move):
        self.move_list.append(move)

        start_x, start_y, end_x, end_y = self.str_to_move(move)

        piece_val = self.board_list[start_x][start_y]
        capt_val = self.board_list[end_x][end_y]

        self.board_list[start_x][start_y] = 0
        self.board_list[end_x][end_y] = piece_val

        en_pass = (piece_val == 1 or piece_val == 7) and start_y != end_y and capt_val == 0
        king_cast = (piece_val == 5 or piece_val == 11) and end_y - start_y == 2
        queen_cast = (piece_val == 5 or piece_val == 11) and start_y - end_y == 2
        promotion = len(move) == 5

        if en_pass:
            #en passant - take pawn
            self.board_list[start_x][end_y] = 0
        elif king_cast:
            #move rook as well
            rook_val = piece_val - 1

            self.board_list[start_x][7] = 0
            self.board_list[start_x][end_y - 1] = rook_val
        elif queen_cast:
            #move rook as well
            rook_val = piece_val - 1

            self.board_list[start_x][0] = 0
            self.board_list[start_x][end_y + 1] = rook_val
        elif promotion:
            #promotion - replace with new piece
            val = graphics_const.BLACK_PIECES.index(move[4])  #promotion is always lowercase

            if piece_val > 6:
                val += 6

            self.board_list[end_x][end_y] = val + 1  #+1 because we are converting from index to piece value

    def update(self, move):
        #update the board after a move has been made
        self.make_move(move)
        self.engine.set_fen(self.start_fen, self.move_list)

    def get_legal_moves(self, start_x, start_y):
        #query the engine to get a list of legal moves as end square coords
        moves = self.engine.get_legal_moves()

        end_squares = []
        for i in moves:
            s_x, s_y, e_x, e_y = self.str_to_move(i)

            if s_x == start_x and s_y == start_y:
                end_squares.append((e_x, e_y))

        return end_squares