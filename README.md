# sklbz-ng

Un serveur Go pour gérer des articles dynamiquement via une API REST et servir des fichiers statiques (JavaScript, CSS). **Utilise maintenant GORM avec SQLite pour le stockage des articles.**

## Fonctionnalités

- **API REST** pour la gestion des articles (CRUD)
- **Base de données SQLite** avec GORM pour le stockage des articles
- **Service de fichiers statiques** pour JavaScript et CSS
- **Filtrage** des articles par auteur, catégorie, tag et visibilité
- **Metadata complète** pour chaque article (auteur, date, visibilité, catégories, tags)
- **Auto-migration** de la base de données au démarrage
- **Arrêt gracieux** du serveur

## Structure du projet

```
sklbz-ng/
├── main.go                 # Point d'entrée principal
├── go.mod                  # Dépendances Go
├── Makefile                # Commandes utiles
├── README.md               # Documentation
├── config/
│   └── database.go         # Configuration de la base de données
├── models/
│   └── article.go          # Modèles Article + repository GORM
├── handlers/
│   └── articles.go         # Handlers API pour les articles
├── server/
│   └── server.go           # Configuration du serveur HTTP
├── static/
│   ├── js/
│   │   └── app.js          # Exemple de fichier JavaScript
│   └── css/
│       └── style.css      # Exemple de fichier CSS
└── sklbz.db               # Base de données SQLite (créée automatiquement)
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
# Avec les paramètres par défaut (port 8080, base de données ./sklbz.db)
./sklbz-ng

# Avec des paramètres personnalisés
./sklbz-ng -port 3000 -db ./my-database.db -static ./my-static

# Avec mTLS activé
./sklbz-ng -mtls -cert certs/server.crt -key certs/server.key -ca-cert certs/ca.crt

# Utiliser le Makefile
make run
```

### Options de la ligne de commande

| Option | Description | Valeur par défaut |
|--------|-------------|------------------|
| `-port` | Port du serveur | `8080` |
| `-db` | Chemin vers la base de données SQLite | `./sklbz.db` |
| `-static` | Répertoire des fichiers statiques | `./static` |
| `-templates` | Répertoire des templates HTML | `./templates` |
| `-mtls` | Activer l'authentification mutuelle TLS | `false` |
| `-cert` | Chemin vers le certificat serveur (PEM) | `""` |
| `-key` | Chemin vers la clé privée serveur (PEM) | `""` |
| `-ca-cert` | Chemin vers le certificat CA pour la vérification client (PEM) | `""` |

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
  "id": 1,
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

## Base de données

### SQLite

Le projet utilise **SQLite** comme base de données par défaut. La base de données est créée automatiquement au premier démarrage.

**Fichier de la base de données :** `./sklbz.db` (ou le chemin spécifié avec `-db`)

### Auto-migration

Au démarrage du serveur, GORM exécute automatiquement les migrations pour créer la table `articles` si elle n'existe pas.

### Schéma de la table articles

```sql
CREATE TABLE articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    author VARCHAR(100) NOT NULL,
    date DATETIME,
    visibility VARCHAR(20) DEFAULT 'public',
    categories TEXT,
    tags TEXT
);
```

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
curl http://localhost:8080/api/articles/1
```

### Mettre à jour un article

```bash
curl -X PUT http://localhost:8080/api/articles/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Introduction à Go - Mise à jour",
    "content": "Go est un langage de programmation moderne et puissant..."
  }'
```

### Supprimer un article

```bash
curl -X DELETE http://localhost:8080/api/articles/1
```

## Sécurité avec mTLS

### Authentification mutuelle TLS

Le serveur supporte l'authentification mutuelle TLS (mTLS) pour sécuriser les opérations d'écriture (POST, PUT, DELETE) sur l'API des articles.

**Fonctionnement :**
- Les requêtes **GET** sont accessibles sans authentification
- Les requêtes **POST, PUT, DELETE** nécessitent un certificat client valide
- Le certificat client doit être signé par une autorité de certification (CA) de confiance

### Génération des certificats

Un script est fourni pour générer des certificats de test :

```bash
# Générer les certificats
./scripts/generate_certs.sh

# Cela crée un répertoire 'certs/' avec :
# - ca.crt, ca.key (Autorité de certification)
# - server.crt, server.key (Certificat serveur)
# - client.crt, client.key (Certificat client)
# - client-combined.pem (Client cert + key pour curl)
```

### Démarrer le serveur avec mTLS

```bash
./sklbz-ng -mtls -cert certs/server.crt -key certs/server.key -ca-cert certs/ca.crt
```

### Tester avec curl

```bash
# Requête GET - pas besoin de certificat
curl http://localhost:8080/api/articles

# Requête POST - nécessite un certificat client
curl -X POST https://localhost:8080/api/articles \
  --cert certs/client-combined.pem \
  --cacert certs/ca.crt \
  -H "Content-Type: application/json" \
  -d '{"title": "Test mTLS", "content": "Contenu test", "author": "Test"}'
```

### Configuration via variables d'environnement

```bash
# Activer mTLS via variable d'environnement
export MTLS_ENABLED=true
./sklbz-ng -cert certs/server.crt -key certs/server.key -ca-cert certs/ca.crt
```

### Sécurité

- **TLS 1.2 minimum** : Le serveur nécessite au moins TLS 1.2
- **Cipher suites sécurisées** : Seules les cipher suites modernes sont autorisées
- **Vérification du certificat** : Le certificat client doit être valide et non expiré
- **Journalisation** : Les requêtes authentifiées sont journalisées avec le CN du certificat

## Développement

### Changer de base de données

Pour utiliser une autre base de données (MySQL, PostgreSQL, etc.), modifiez le fichier `config/database.go` :

```go
// Exemple pour MySQL
import "gorm.io/driver/mysql"

func InitDB(dsn string) *gorm.DB {
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    // ...
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
