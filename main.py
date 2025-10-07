import typer
from rich import print

from parser.parse import get_ast


def extract_text_from_test():
    """Extract text from test.py using hardcoded path"""
    test_file_path = "test.py"

    try:
        with open(test_file_path, "r", encoding="utf-8") as file:
            text = file.read()
        return text
    except FileNotFoundError:
        print(f"Error: File '{test_file_path}' not found")
        return None
    except Exception as e:
        print(f"Error reading file: {e}")
        return None


def main():
    text = extract_text_from_test()
    if text:
        get_ast(text)


if __name__ == "__main__":
    typer.run(main)
