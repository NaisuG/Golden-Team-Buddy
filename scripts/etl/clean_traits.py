import json

INPUT_FILE = "raw_traits.json"
OUTPUT_FILE = "traits.json"

FIELDS_TO_KEEP = ["key", "ingameKey", "name", "imageUrl", "type", "styles", "stats"]


def clean_trait(raw):
    return {field: raw.get(field) for field in FIELDS_TO_KEEP}


def main():
    with open(INPUT_FILE, "r", encoding="utf-8") as f:
        data = json.load(f)

    playable = [t for t in data["traits"] if not t.get("isHidden")]
    cleaned = [clean_trait(t) for t in playable]

    output = {"season": data["season"], "traits": cleaned}

    with open(OUTPUT_FILE, "w", encoding="utf-8") as f:
        json.dump(output, f, indent=2, ensure_ascii=False)

    print(f"Total en raw: {len(data['traits'])}")
    print(f"Despues de filtrar isHidden: {len(playable)}")
    print(f"Guardado en {OUTPUT_FILE}")


if __name__ == "__main__":
    main()