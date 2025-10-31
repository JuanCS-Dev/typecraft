#!/usr/bin/env python3
"""
Fetch Google Fonts metadata and create a curated database for font suggestions.
"""

import json
import requests
from typing import List, Dict, Any

# Top curated fonts for book typography
CURATED_FONTS = {
    "serif_body": [
        "Crimson Text",
        "Merriweather",
        "Lora",
        "Libre Baskerville",
        "Source Serif Pro",
        "Playfair Display",
        "EB Garamond",
        "Cormorant Garamond",
        "Spectral",
        "Gelasio",
    ],
    "sans_body": [
        "Source Sans Pro",
        "Open Sans",
        "Roboto",
        "Lato",
        "Montserrat",
        "Raleway",
        "PT Sans",
        "Work Sans",
        "Nunito",
        "Inter",
    ],
    "serif_heading": [
        "Playfair Display",
        "Merriweather",
        "Crimson Text",
        "Libre Baskerville",
        "Cormorant Garamond",
        "Cinzel",
        "Cardo",
        "Literata",
        "Unna",
        "Vollkorn",
    ],
    "sans_heading": [
        "Montserrat",
        "Oswald",
        "Raleway",
        "Roboto Condensed",
        "Archivo Black",
        "Anton",
        "Bebas Neue",
        "Barlow Condensed",
        "Saira Condensed",
        "Fjalla One",
    ],
    "monospace": [
        "Fira Code",
        "Source Code Pro",
        "Roboto Mono",
        "JetBrains Mono",
        "IBM Plex Mono",
        "Inconsolata",
        "Courier Prime",
        "Space Mono",
        "Anonymous Pro",
        "Ubuntu Mono",
    ],
}

# Genre-based font pairings
GENRE_PAIRINGS = {
    "fiction": [
        {"body": "Crimson Text", "heading": "Playfair Display", "mood": "classic", "rationale": "Elegant serif pair for literary fiction"},
        {"body": "Merriweather", "heading": "Montserrat", "mood": "modern", "rationale": "Contemporary serif/sans mix"},
        {"body": "Lora", "heading": "Raleway", "mood": "clean", "rationale": "Readable and approachable"},
    ],
    "mystery": [
        {"body": "Crimson Text", "heading": "Playfair Display", "mood": "dark", "rationale": "Classic noir aesthetic"},
        {"body": "Libre Baskerville", "heading": "Oswald", "mood": "tense", "rationale": "Sharp contrasts for suspense"},
        {"body": "EB Garamond", "heading": "Cinzel", "mood": "sophisticated", "rationale": "Refined detective style"},
    ],
    "romance": [
        {"body": "Lora", "heading": "Playfair Display", "mood": "romantic", "rationale": "Soft and elegant"},
        {"body": "Crimson Text", "heading": "Cormorant Garamond", "mood": "delicate", "rationale": "Flowing and graceful"},
        {"body": "Spectral", "heading": "Montserrat", "mood": "contemporary", "rationale": "Modern romance"},
    ],
    "scifi": [
        {"body": "Source Sans Pro", "heading": "Roboto Condensed", "mood": "futuristic", "rationale": "Clean tech aesthetic"},
        {"body": "Inter", "heading": "Oswald", "mood": "modern", "rationale": "Sharp and forward-looking"},
        {"body": "Work Sans", "heading": "Barlow Condensed", "mood": "minimalist", "rationale": "Streamlined future"},
    ],
    "fantasy": [
        {"body": "Merriweather", "heading": "Cinzel", "mood": "epic", "rationale": "Grand and mythical"},
        {"body": "Crimson Text", "heading": "Playfair Display", "mood": "classic", "rationale": "Timeless storytelling"},
        {"body": "Libre Baskerville", "heading": "Unna", "mood": "mystical", "rationale": "Enchanted aesthetic"},
    ],
    "technical": [
        {"body": "Source Sans Pro", "heading": "Roboto", "mono": "Fira Code", "mood": "professional", "rationale": "Clear technical documentation"},
        {"body": "Inter", "heading": "Montserrat", "mono": "Source Code Pro", "mood": "modern", "rationale": "Contemporary tech writing"},
        {"body": "Open Sans", "heading": "Raleway", "mono": "JetBrains Mono", "mood": "accessible", "rationale": "User-friendly technical content"},
    ],
    "academic": [
        {"body": "Source Serif Pro", "heading": "Roboto Slab", "mono": "Inconsolata", "mood": "scholarly", "rationale": "Academic professionalism"},
        {"body": "Crimson Text", "heading": "Libre Baskerville", "mood": "traditional", "rationale": "Classical academic style"},
        {"body": "EB Garamond", "heading": "Montserrat", "mood": "refined", "rationale": "Sophisticated research"},
    ],
    "business": [
        {"body": "Lato", "heading": "Montserrat", "mood": "corporate", "rationale": "Professional business style"},
        {"body": "Source Sans Pro", "heading": "Oswald", "mood": "confident", "rationale": "Strong business voice"},
        {"body": "Roboto", "heading": "Raleway", "mood": "modern", "rationale": "Contemporary business"},
    ],
}

def create_font_database() -> Dict[str, Any]:
    """Create the font database with curated fonts and pairings."""
    
    database = {
        "version": "1.0",
        "last_updated": "2024-10-31",
        "curated_fonts": CURATED_FONTS,
        "genre_pairings": GENRE_PAIRINGS,
        "font_metadata": {},
    }
    
    # Add metadata for each curated font
    all_fonts = set()
    for category in CURATED_FONTS.values():
        all_fonts.update(category)
    
    for font in all_fonts:
        database["font_metadata"][font] = {
            "name": font,
            "category": get_font_category(font),
            "supports_latin": True,
            "variable": False,
            "google_fonts": True,
        }
    
    return database

def get_font_category(font_name: str) -> str:
    """Determine the category of a font."""
    serif_keywords = ["garamond", "baskerville", "serif", "crimson", "merriweather", 
                      "lora", "playfair", "spectral", "gelasio", "cinzel", "cardo",
                      "literata", "unna", "vollkorn"]
    
    mono_keywords = ["mono", "code", "courier"]
    
    font_lower = font_name.lower()
    
    if any(k in font_lower for k in mono_keywords):
        return "monospace"
    elif any(k in font_lower for k in serif_keywords):
        return "serif"
    else:
        return "sans-serif"

def main():
    """Generate and save the font database."""
    print("🎨 Creating Google Fonts Database...")
    
    database = create_font_database()
    
    # Save to JSON
    output_path = "../pkg/design/google_fonts_db.json"
    with open(output_path, 'w', encoding='utf-8') as f:
        json.dump(database, f, indent=2, ensure_ascii=False)
    
    print(f"✅ Database created: {output_path}")
    print(f"   - {len(database['curated_fonts'])} font categories")
    print(f"   - {len(database['font_metadata'])} curated fonts")
    print(f"   - {len(database['genre_pairings'])} genre pairings")
    print(f"   - Total pairings: {sum(len(p) for p in database['genre_pairings'].values())}")

if __name__ == "__main__":
    main()
