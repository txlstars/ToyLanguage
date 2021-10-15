#!/usr/bin/env python3

import argparse
import subprocess

args = None


def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--glp", help="generate lexer and parser source code by g4", action="store_true"
    )
    global args
    args = parser.parse_args()


antlr4 = "java -jar ./tools/antlr-4.9.2-complete.jar"

def generate_lexer_and_parser():
    subprocess.run(
        "{} -Xexact-output-dir -o ./parser -package parser -Dlanguage=Go ./parser/*.g4".format(
            antlr4
        ),
        shell=True,
        check=True,
        capture_output=True,
    )


def main():
    parse_args()

    if args.glp:
        generate_lexer_and_parser()

if __name__ == "__main__":
    main()
