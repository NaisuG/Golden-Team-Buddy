import json

INPUT_FILE = "raw_items.json"
OUTPUT_FILE = "items.json"

FIELDS_TO_KEEP = ["key", "ingameKey", "name", "desc", "fromDesc", "imageUrl", "compositions", "affectedTraitKey"]


def categorize(raw):
    if raw.get("isFromItem"):
        return "component"
    if raw.get("isRadiant"):
        return "radiant"
    if raw.get("isEmblem"):
        return "emblem"
    if raw.get("isArtifact"):
        return "artifact"
    if raw.get("isNormal"):
        return "completed"
    if raw.get("isTrait") or raw.get("isPsionic") or raw.get("isAnimaSquad") or raw.get("isCommandMod"):
        return "special"
    return "other"


def clean_item(raw):
    cleaned = {field: raw.get(field) for field in FIELDS_TO_KEEP}
    cleaned["category"] = categorize(raw)
    return cleaned


def main():
    with open(INPUT_FILE, "r", encoding="utf-8") as f:
        data = json.load(f)

    playable = [i for i in data["items"] if not i.get("isHidden")]
    cleaned = [clean_item(i) for i in playable]

    output = {"season": data["season"], "items": cleaned}

    with open(OUTPUT_FILE, "w", encoding="utf-8") as f:
        json.dump(output, f, indent=2, ensure_ascii=False)

    print(f"Total en raw: {len(data['items'])}")
    print(f"Despues de filtrar isHidden: {len(playable)}")
    print(f"Guardado en {OUTPUT_FILE}")


if __name__ == "__main__":
    main()