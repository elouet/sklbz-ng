# sklbz-ng

Un serveur Go pour gérer des articles dynamiquement via une API REST et servir des fichiers statiques (JavaScript, CSS).

## Fonctionnalités

- **API REST** pour la gestion des articles (CRUD)
- **Stockage JSON** sur disque pour les articles
- **Service de fichiers statiques** pour JavaScript et CSS
- **Filtrage** des articles par auteur, catégorie, tag et visibilité
- **Metadata complète** pour chaque article (auteur, date, visibilité, catégories, tags)

## Structure du projet

```
sklbz-ng/
├── main.go                 # Point d'entrée principal
├── go.mod                  # Dépendances Go
├── Makefile                # Commandes utiles
├── README.md               # Documentation
├── models/
│   └── article.go          # Modèles et stockage des articles
├── handlers/
│   └── articles.go         # Handlers API pour les articles
├── server/
│   └── server.go           # Configuration du serveur HTTP
├── static/
│   ├── js/
│   │   └── app.js          # Exemple de fichier JavaScript
│   └── css/
│       └── style.css      # Exemple de fichier CSS
└── data/                   # Stockage des articles (fichiers JSON)
```

## Installation

### Prérequis

- Go 1.21 ou supérieur
- Git

### Installation

```bash
# Cloner le dépôt
git clone https://github.com/elouet/sklbz-ng.git
cd sklbz-ng

# Télécharger les dépendances
go mod download

# Construire l'application
go build -o sklbz-ng
```

## Utilisation

### Démarrer le serveur

```bash
# Avec les paramètres par défaut (port 8080)
./sklbz-ng

# Avec des paramètres personnalisés
./sklbz-ng -port 3000 -data ./my-data -static ./my-static

# Utiliser le Makefile
make run
```

### Options de la ligne de commande

| Option | Description | Valeur par défaut |
|--------|-------------|------------------|
| `-port` | Port du serveur | `8080` |
| `-data` | Répertoire de stockage des articles | `./data` |
| `-static` | Répertoire des fichiers statiques | `./static` |

## API Endpoints

### Articles

| Méthode | Endpoint | Description |
|---------|----------|-------------|
| GET | `/api/articles` | Lister tous les articles (avec filtrage possible) |
| POST | `/api/articles` | Créer un nouvel article |
| GET | `/api/articles/{id}` | Récupérer un article par ID |
| PUT | `/api/articles/{id}` | Mettre à jour un article |
| DELETE | `/api/articles/{id}` | Supprimer un article |

### Filtrage des articles

Vous pouvez filtrer les articles en utilisant les paramètres de requête :

- `?author=<nom>` - Filtrer par auteur
- `?category=<catégorie>` - Filtrer par catégorie
- `?tag=<tag>` - Filtrer par tag
- `?visibility=<visibilité>` - Filtrer par visibilité (public, private, draft)

Exemple : `/api/articles?author=John&category=technology&tag=go`

### Fichiers statiques

| Endpoint | Description |
|----------|-------------|
| `/static/js/` | Accéder aux fichiers JavaScript |
| `/static/css/` | Accéder aux fichiers CSS |
| `/static/` | Accéder aux autres fichiers statiques |

### Autres endpoints

| Méthode | Endpoint | Description |
|---------|----------|-------------|
| GET | `/` | Message de bienvenue avec les endpoints disponibles |
| GET | `/health` | Vérification de l'état du serveur |

## Format des articles

### Requête de création/mise à jour

```json
{
  "title": "Mon premier article",
  "content": "Ceci est le contenu de mon article.",
  "author": "Jean Dupont",
  "visibility": "public",
  "categories": ["technologie", "web"],
  "tags": ["go", "api", "blog"]
}
```

### Réponse d'article

```json
{
  "id": "article-1234567890",
  "title": "Mon premier article",
  "content": "Ceci est le contenu de mon article.",
  "author": "Jean Dupont",
  "date": "2024-01-15T10:30:00Z",
  "visibility": "public",
  "categories": ["technologie", "web"],
  "tags": ["go", "api", "blog"],
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Visibilité

Les articles peuvent avoir les niveaux de visibilité suivants :
- `public` - Visible par tous
- `private` - Visible uniquement par l'auteur
- `draft` - Brouillon (non publié)

## Exemples d'utilisation

### Créer un article

```bash
curl -X POST http://localhost:8080/api/articles \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Introduction à Go",
    "content": "Go est un langage de programmation moderne...",
    "author": "Alice",
    "visibility": "public",
    "categories": ["programmation", "go"],
    "tags": ["tutoriel", "débutant"]
  }'
```

### Lister tous les articles

```bash
curl http://localhost:8080/api/articles
```

### Filtrer les articles

```bash
# Par auteur
curl http://localhost:8080/api/articles?author=Alice

# Par catégorie
curl http://localhost:8080/api/articles?category=programmation

# Par tag
curl http://localhost:8080/api/articles?tag=tutoriel

# Combinaison
curl http://localhost:8080/api/articles?author=Alice&category=programmation
```

### Récupérer un article spécifique

```bash
curl http://localhost:8080/api/articles/article-1234567890
```

### Mettre à jour un article

```bash
curl -X PUT http://localhost:8080/api/articles/article-1234567890 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Introduction à Go - Mise à jour",
    "content": "Go est un langage de programmation moderne et puissant..."
  }'
```

### Supprimer un article

```bash
curl -X DELETE http://localhost:8080/api/articles/article-1234567890
```

## Développement

### Structure des fichiers

Les articles sont stockés sous forme de fichiers JSON dans le répertoire `data/` :

```json
{
  "id": "article-1234567890",
  "title": "Mon premier article",
  "content": "Ceci est le contenu...",
  "author": "Jean Dupont",
  "date": "2024-01-15T10:30:00Z",
  "visibility": "public",
  "categories": ["technologie", "web"],
  "tags": ["go", "api"],
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Ajouter des fichiers statiques

Placez vos fichiers JavaScript dans `static/js/` et vos fichiers CSS dans `static/css/`.

Exemple :
- `static/js/app.js` → accessible via `/static/js/app.js`
- `static/css/style.css` → accessible via `/static/css/style.css`

## Commandes Makefile

| Commande | Description |
|----------|-------------|
| `make build` | Construire l'application |
| `make run` | Construire et exécuter |
| `make clean` | Nettoyer les artefacts de construction |
| `make test` | Exécuter les tests |
| `make deps` | Télécharger les dépendances |
| `make run-dev` | Exécuter en mode développement |
| `make build-prod` | Construire pour la production |
| `make help` | Afficher l'aide |

## Contribution

Les contributions sont les bienvenues ! Veuillez ouvrir une Pull Request avec vos modifications.

## Licence

MIT
