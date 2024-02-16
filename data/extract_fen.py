import chess.pgn
from io import StringIO


def read_file():
    with open("lichess_db_standard_rated_2013-01.txt", "r") as file:
        data = file.read()

    return data.split("\n")


def write_file(string):
    with open("fen_data.txt", "w") as file:
        file.write(string)


def extract_games(data):
    games = []
    headers = []
    for i in data:
        if i == "":
            continue

        if i[0] == "[":
            headers.append(i)
        else:
            head_str = "\n".join(headers)
            headers = []

            game_str = f"{head_str}\n\n{i}"
            games.append(game_str)

    return games


def get_fens(pgn_game):
    game = chess.pgn.read_game(StringIO(pgn_game))
    board = game.board()
    moves = list(game.mainline_moves())

    fens = []
    for i in range(0, len(moves) - 1, 2):
        board.push(moves[i])
        board.push(moves[i + 1])

        fens.append(board.fen())

    return fens


def get_all_fen(num_games):
    data = read_file()
    games = extract_games(data)[:num_games]

    fens = set()  #use set so duplicates not included
    for i, x in enumerate(games):
        for j in get_fens(x):
            fens.add(j)

        done_frac = i / len(games)
        num_hash = int(75 * done_frac)

        print(f"Extracting game fens: |{'#' * num_hash}{'.' * (75 - num_hash)}| {done_frac * 100 :.2f}%", end="\r")

    print("\nDone!")

    return fens


def main():
    fens = get_all_fen(50000)

    write_file("\n".join(fens))


main()