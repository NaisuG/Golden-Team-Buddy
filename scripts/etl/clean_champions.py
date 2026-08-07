import json

INPUT_FILE = "raw_champions.json"
OUTPUT_FILE = "champions.json"

FIELDS_TO_KEEP = ["key", "ingameKey", "name", "imageUrl", "traits", "role"]


def clean_champion(raw):
    cleaned = {field: raw.get(field) for field in FIELDS_TO_KEEP}
    cleaned["cost"] = raw["cost"][0]
    return cleaned


def main():
    with open(INPUT_FILE, "r", encoding="utf-8") as f:
        data = json.load(f)

    playable = [c for c in data["champions"] if not c.get("isHidden")]
    cleaned = [clean_champion(c) for c in playable]

    output = {"season": data["season"], "champions": cleaned}

    with open(OUTPUT_FILE, "w", encoding="utf-8") as f:
        json.dump(output, f, indent=2, ensure_ascii=False)

    print(f"Total en raw: {len(data['champions'])}")
    print(f"Despues de filtrar isHidden: {len(playable)}")
    print(f"Guardado en {OUTPUT_FILE}")


if __name__ == "__main__":
    main()