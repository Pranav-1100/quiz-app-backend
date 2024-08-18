package database

import (
    "database/sql"
    "log"
    "strings"
    _ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }

    if err := runMigrations(db); err != nil {
        return nil, err
    }

    if err := InsertMockData(db); err != nil {
        log.Printf("Error inserting mock data: %v", err)
        // Don't return here, as the table might already have data
    }

    return db, nil
}

func runMigrations(db *sql.DB) error {
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS questions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            level_id INTEGER,
            question TEXT NOT NULL,
            options TEXT NOT NULL,
            correct_answer TEXT NOT NULL,
            difficulty INTEGER NOT NULL,
            image_url TEXT,
            explanation TEXT
        );

        CREATE TABLE IF NOT EXISTS user_progress (
            user_id INTEGER,
            level_id INTEGER,
            completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            PRIMARY KEY (user_id, level_id)
        );

        CREATE TABLE IF NOT EXISTS levels (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            difficulty INTEGER NOT NULL,
            unlock_threshold INTEGER NOT NULL,
            description TEXT
        );

        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT UNIQUE NOT NULL,
            coins INTEGER DEFAULT 0,
            current_level INTEGER DEFAULT 1,
            streak INTEGER DEFAULT 0
        );

        CREATE TABLE IF NOT EXISTS user_achievements (
            user_id INTEGER,
            achievement_id INTEGER,
            unlocked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            PRIMARY KEY (user_id, achievement_id)
        );

        CREATE TABLE IF NOT EXISTS leaderboard (
            user_id INTEGER,
            score INTEGER,
            timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            PRIMARY KEY (user_id, timestamp)
        );
    `)

    if err != nil {
        // If the error is about the column already existing, we can ignore it
        if !strings.Contains(err.Error(), "duplicate column name") {
            return err
        }
    }

    return nil
}

func InsertMockData(db *sql.DB) error {
    _, err := db.Exec(`
    -- Insert mock niches with difficulty, unlock_threshold, and description
    INSERT OR IGNORE INTO levels (name, difficulty, unlock_threshold, description) VALUES 
    ('General Knowledge', 1, 0, 'Test your knowledge on a variety of topics'),
    ('Science', 1, 10, 'Explore the wonders of the natural world'),
    ('History', 1, 20, 'Journey through time and learn about past events'),
    ('Geography', 2, 30, 'Discover the world''s places and features'),
    ('Literature', 2, 40, 'Dive into the world of books and authors'),
    ('Movies', 2, 50, 'Test your knowledge of cinema and films'),
    ('Music', 3, 60, 'Explore the world of melodies and rhythms'),
    ('Sports', 3, 70, 'Challenge yourself with questions about various sports'),
    ('Technology', 3, 80, 'Navigate through the world of gadgets and innovations'),
    ('Art', 3, 90, 'Immerse yourself in the world of creativity and expression');

-- 1. General Knowledge
INSERT INTO questions (level_id, question, options, correct_answer, difficulty, image_url, explanation) VALUES
(1, 'What is the capital of France?', '["London", "Berlin", "Paris", "Madrid"]', 'Paris', 1, 'https://example.com/paris.jpg', 'Paris is the capital and most populous city of France.'),
(1, 'Which planet is known as the Red Planet?', '["Venus", "Mars", "Jupiter", "Saturn"]', 'Mars', 1, 'https://example.com/mars.jpg', 'Mars is often called the Red Planet due to its reddish appearance.'),
(1, 'What is the largest mammal in the world?', '["Elephant", "Blue Whale", "Giraffe", "Hippopotamus"]', 'Blue Whale', 1, 'https://example.com/blue_whale.jpg', 'The blue whale is the largest animal known to have ever existed.'),
(1, 'Who wrote "Romeo and Juliet"?', '["Charles Dickens", "William Shakespeare", "Jane Austen", "Mark Twain"]', 'William Shakespeare', 1, 'https://example.com/shakespeare.jpg', 'William Shakespeare wrote "Romeo and Juliet" in the late 16th century.'),
(1, 'What is the chemical symbol for gold?', '["Au", "Ag", "Fe", "Cu"]', 'Au', 2, 'https://example.com/gold.jpg', 'Au is the chemical symbol for gold, derived from its Latin name "aurum".'),
(1, 'Which country has the largest population in the world?', '["India", "United States", "China", "Russia"]', 'China', 2, 'https://example.com/china.jpg', 'As of 2024, China has the largest population in the world, although India is a close second.'),
(1, 'What is the speed of light?', '["299,792 km/s", "300,000 km/s", "150,000 km/s", "200,000 km/s"]', '299,792 km/s', 2, 'https://example.com/speed_of_light.jpg', 'The speed of light in vacuum is exactly 299,792 kilometers per second.'),
(1, 'Who developed the theory of relativity?', '["Isaac Newton", "Albert Einstein", "Stephen Hawking", "Niels Bohr"]', 'Albert Einstein', 3, 'https://example.com/einstein.jpg', 'Albert Einstein developed both the special and general theories of relativity.'),
(1, 'What is the Fibonacci sequence?', '["Sum of two preceding numbers", "Sequence of prime numbers", "Geometric sequence", "Perfect squares"]', 'Sum of two preceding numbers', 3, 'https://example.com/fibonacci.jpg', 'The Fibonacci sequence is a series where each number is the sum of the two preceding ones.'),
(1, 'What is the half-life of Carbon-14?', '["2,730 years", "5,730 years", "7,730 years", "10,730 years"]', '5,730 years', 3, 'https://example.com/carbon14.jpg', 'The half-life of Carbon-14 is approximately 5,730 years, used for dating organic materials.');

-- 2. Science
INSERT INTO questions (level_id, question, options, correct_answer, difficulty, image_url, explanation) VALUES
(2, 'What is the chemical formula for water?', '["H2O", "CO2", "NaCl", "CH4"]', 'H2O', 1, 'https://example.com/water.jpg', 'Water is composed of two hydrogen atoms and one oxygen atom.'),
(2, 'What is the largest organ in the human body?', '["Heart", "Brain", "Liver", "Skin"]', 'Skin', 1, 'https://example.com/skin.jpg', 'The skin is the largest organ of the human body by weight and surface area.'),
(2, 'What is the process by which plants make their own food?', '["Photosynthesis", "Respiration", "Fermentation", "Digestion"]', 'Photosynthesis', 1, 'https://example.com/photosynthesis.jpg', 'Photosynthesis is the process by which plants use sunlight to produce energy from carbon dioxide and water.'),
(2, 'What is the hardest natural substance on Earth?', '["Gold", "Iron", "Diamond", "Titanium"]', 'Diamond', 1, 'https://example.com/diamond.jpg', 'Diamond is the hardest known natural material on Earth.'),
(2, 'What is the unit of electrical resistance?', '["Volt", "Ampere", "Watt", "Ohm"]', 'Ohm', 2, 'https://example.com/ohm.jpg', 'The ohm is the SI unit of electrical resistance.'),
(2, 'What is the process of cell division called?', '["Mitosis", "Meiosis", "Photosynthesis", "Respiration"]', 'Mitosis', 2, 'https://example.com/mitosis.jpg', 'Mitosis is the process of cell division which results in two identical daughter cells.'),
(2, 'What is the pH of a neutral solution?', '["0", "7", "14", "10"]', '7', 2, 'https://example.com/ph_scale.jpg', 'On the pH scale, 7 is considered neutral, with lower values being acidic and higher values being basic.'),
(2, 'What force keeps planets in orbit around the Sun?', '["Electromagnetic", "Strong nuclear", "Weak nuclear", "Gravitational"]', 'Gravitational', 3, 'https://example.com/gravity.jpg', 'The gravitational force, as described by Newton''s law of universal gravitation, keeps planets in orbit.'),
(2, 'What is the function of mitochondria in a cell?', '["Protein synthesis", "Energy production", "Storage", "Cell division"]', 'Energy production', 3, 'https://example.com/mitochondria.jpg', 'Mitochondria are often called the powerhouses of the cell because they produce most of the cell''s ATP.'),
(2, 'What is Heisenberg''s Uncertainty Principle?', '["Energy of a photon", "Conservation of energy", "Position and momentum uncertainty", "Wave-particle duality"]', 'Position and momentum uncertainty', 3, 'https://example.com/uncertainty_principle.jpg', 'It states that the position and momentum of a particle cannot be simultaneously measured with arbitrary precision.');

-- 3. History
INSERT INTO questions (level_id, question, options, correct_answer, difficulty, image_url, explanation) VALUES
(3, 'In which year did World War II end?', '["1943", "1945", "1947", "1950"]', '1945', 1, 'https://example.com/ww2_end.jpg', 'World War II ended in 1945 with the surrender of Germany in May and Japan in August.'),
(3, 'Who was the first President of the United States?', '["Thomas Jefferson", "John Adams", "George Washington", "Benjamin Franklin"]', 'George Washington', 1, 'https://example.com/washington.jpg', 'George Washington was the first President of the United States, serving from 1789 to 1797.'),
(3, 'What ancient wonder was located in Alexandria, Egypt?', '["Hanging Gardens", "Colossus of Rhodes", "Lighthouse", "Temple of Artemis"]', 'Lighthouse', 1, 'https://example.com/alexandria_lighthouse.jpg', 'The Lighthouse of Alexandria was one of the Seven Wonders of the Ancient World.'),
(3, 'Who painted the Mona Lisa?', '["Vincent van Gogh", "Leonardo da Vinci", "Pablo Picasso", "Michelangelo"]', 'Leonardo da Vinci', 1, 'https://example.com/mona_lisa.jpg', 'The Mona Lisa was painted by Italian Renaissance artist Leonardo da Vinci.'),
(3, 'What was the name of the first successful English settlement in America?', '["Plymouth", "Roanoke", "Jamestown", "New Amsterdam"]', 'Jamestown', 2, 'https://example.com/jamestown.jpg', 'Jamestown, founded in 1607, was the first permanent English settlement in North America.'),
(3, 'Which empire was ruled by Genghis Khan?', '["Roman Empire", "Ottoman Empire", "Mongol Empire", "Byzantine Empire"]', 'Mongol Empire', 2, 'https://example.com/genghis_khan.jpg', 'Genghis Khan founded and ruled the Mongol Empire, which became the largest contiguous land empire in history.'),
(3, 'What event marked the start of World War I?', '["Russian Revolution", "Assassination of Archduke Franz Ferdinand", "German invasion of Poland", "Bombing of Pearl Harbor"]', 'Assassination of Archduke Franz Ferdinand', 2, 'https://example.com/franz_ferdinand.jpg', 'The assassination of Archduke Franz Ferdinand of Austria in 1914 is considered the immediate trigger of World War I.'),
(3, 'What was the significance of the Rosetta Stone?', '["Ancient map", "Deciphering hieroglyphics", "Religious artifact", "Royal decree"]', 'Deciphering hieroglyphics', 3, 'https://example.com/rosetta_stone.jpg', 'The Rosetta Stone provided the key to understanding Egyptian hieroglyphics.'),
(3, 'Who was the last Emperor of Russia?', '["Peter the Great", "Ivan the Terrible", "Nicholas II", "Alexander III"]', 'Nicholas II', 3, 'https://example.com/nicholas_ii.jpg', 'Nicholas II was the last Emperor of Russia, ruling from 1894 until his forced abdication in 1917.'),
(3, 'What was the main cause of the French Revolution?', '["Foreign invasion", "Religious conflict", "Social and economic inequality", "Natural disaster"]', 'Social and economic inequality', 3, 'https://example.com/french_revolution.jpg', 'The French Revolution was primarily caused by social and economic inequalities among the classes of French society.');

-- 4. Geography
INSERT INTO questions (level_id, question, options, correct_answer, difficulty, image_url, explanation) VALUES
(4, 'What is the largest continent by land area?', '["North America", "Africa", "Europe", "Asia"]', 'Asia', 1, 'https://example.com/asia.jpg', 'Asia is the largest continent, covering approximately 30% of Earth''s total land area.'),
(4, 'Which river is the longest in the world?', '["Amazon", "Nile", "Yangtze", "Mississippi"]', 'Nile', 1, 'https://example.com/nile.jpg', 'The Nile is the longest river in the world, stretching about 6,650 kilometers (4,132 miles).'),
(4, 'What is the capital of Japan?', '["Seoul", "Beijing", "Tokyo", "Bangkok"]', 'Tokyo', 1, 'https://example.com/tokyo.jpg', 'Tokyo is the capital and most populous city of Japan.'),
(4, 'Which country is home to the Great Barrier Reef?', '["Brazil", "Australia", "Indonesia", "Philippines"]', 'Australia', 1, 'https://example.com/great_barrier_reef.jpg', 'The Great Barrier Reef is located off the coast of Queensland in northeastern Australia.'),
(4, 'What is the highest mountain in North America?', '["Mount McKinley", "Mount Logan", "Mount Rainier", "Mount Whitney"]', 'Mount McKinley', 2, 'https://example.com/mount_mckinley.jpg', 'Mount McKinley, also known as Denali, is the highest peak in North America at 20,310 feet (6,190 meters).'),
(4, 'Which country has the most time zones?', '["Russia", "United States", "France", "Australia"]', 'France', 2, 'https://example.com/france_time_zones.jpg', 'France has 12 time zones, the most of any country, due to its overseas territories.'),
(4, 'What is the driest place on Earth?', '["Sahara Desert", "Atacama Desert", "Death Valley", "Antarctic Dry Valleys"]', 'Antarctic Dry Valleys', 2, 'https://example.com/antarctic_dry_valleys.jpg', 'The McMurdo Dry Valleys in Antarctica are the driest place on Earth, with extremely low humidity and almost no snow or ice cover.'),
(4, 'What is the deepest point in Earth''s oceans?', '["Mariana Trench", "Tonga Trench", "Philippine Trench", "Puerto Rico Trench"]', 'Mariana Trench', 3, 'https://example.com/mariana_trench.jpg', 'The Challenger Deep in the Mariana Trench is the deepest known point in Earth''s oceans, at about 36,200 feet (11,000 meters) deep.'),
(4, 'Which African country has the most pyramids?', '["Egypt", "Sudan", "Libya", "Algeria"]', 'Sudan', 3, 'https://example.com/sudan_pyramids.jpg', 'Sudan has more pyramids than Egypt, with over 200 pyramids compared to Egypt''s 138.'),
(4, 'What is the only continent that lies in all four hemispheres?', '["Asia", "Africa", "South America", "Australia"]', 'Africa', 3, 'https://example.com/africa_hemispheres.jpg', 'Africa is the only continent to extend into all four hemispheres: Northern, Southern, Eastern, and Western.');

-- 5. Literature
INSERT INTO questions (level_id, question, options, correct_answer, difficulty, image_url, explanation) VALUES
(5, 'Who wrote "To Kill a Mockingbird"?', '["Ernest Hemingway", "Harper Lee", "F. Scott Fitzgerald", "John Steinbeck"]', 'Harper Lee', 1, 'https://example.com/harper_lee.jpg', 'Harper Lee wrote the Pulitzer Prize-winning novel "To Kill a Mockingbird", published in 1960.'),
(5, 'What is the name of the hobbit in "The Lord of the Rings"?', '["Bilbo", "Frodo", "Sam", "Pippin"]', 'Frodo', 1, 'https://example.com/frodo.jpg', 'Frodo Baggins is the primary hobbit protagonist in J.R.R. Tolkien''s "The Lord of the Rings".'),
(5, 'Who wrote "1984"?', '["Aldous Huxley", "Ray Bradbury", "George Orwell", "Philip K. Dick"]', 'George Orwell', 1, 'https://example.com/george_orwell.jpg', 'George Orwell wrote the dystopian novel "1984", published in 1949.'),
(5, 'Which Shakespeare play features the character Juliet?', '["Hamlet", "Macbeth", "Romeo and Juliet", "Othello"]', 'Romeo and Juliet', 1, 'https://example.com/romeo_and_juliet.jpg', 'Juliet is one of the main characters in William Shakespeare''s play "Romeo and Juliet".'),
(5, 'Who wrote "The Catcher in the Rye"?', '["J.D. Salinger", "Jack Kerouac", "William Faulkner", "Kurt Vonnegut"]', 'J.D. Salinger', 2, 'https://example.com/jd_salinger.jpg', 'J.D. Salinger wrote "The Catcher in the Rye", published in 1951.'),
(5, 'What is the name of the fictional country where George Orwell''s "Animal Farm" is set?', '["Manor Farm", "Home Farm", "State Farm", "Plantation Farm"]', 'Manor Farm', 2, 'https://example.com/animal_farm.jpg', 'The story of "Animal Farm" takes place on Manor Farm, which is later renamed Animal Farm by the animals.'),
(5, 'Who wrote "One Hundred Years of Solitude"?', '["Pablo Neruda", "Gabriel García Márquez", "Jorge Luis Borges", "Isabel Allende"]', 'Gabriel García Márquez', 2, 'https://example.com/gabriel_garcia_marquez.jpg', 'Gabriel García Márquez wrote "One Hundred Years of Solitude", a landmark of magical realism.'),
(5, 'What is the first line of Jane Austen''s "Pride and Prejudice"?', '["It was the best of times, it was the worst of times", "Happy families are all alike; every unhappy family is unhappy in its own way", "It is a truth universally acknowledged, that a single man in possession of a good fortune, must be in want of a wife", "The past is a foreign country; they do things differently there"]', 'It is a truth universally acknowledged, that a single man in possession of a good fortune, must be in want of a wife', 3, 'https://example.com/pride_and_prejudice.jpg', 'This famous opening line sets the tone for Jane Austen''s novel about marriage and social status.'),
(5, 'Who wrote "Waiting for Godot"?', '["Samuel Beckett", "Eugene Ionesco", "Harold Pinter", "Tom Stoppard"]', 'Samuel Beckett', 3, 'https://example.com/samuel_beckett.jpg', 'Samuel Beckett wrote "Waiting for Godot", a seminal work in the Theatre of the Absurd.'),
(5, 'What is the significance of the number 42 in "The Hitchhiker''s Guide to the Galaxy"?', '["The age of the protagonist", "The number of planets visited", "The answer to life, the universe, and everything", "The code to start the spaceship"]', 'The answer to life, the universe, and everything', 3, 'https://example.com/hitchhikers_guide.jpg', 'In Douglas Adams'' series, 42 is famously the "Answer to the Ultimate Question of Life, the Universe, and Everything".');

-- 6. Movies
INSERT INTO questions (level_id, question, options, correct_answer, difficulty, image_url, explanation) VALUES
(6, 'Who directed the movie "Jaws"?', '["Martin Scorsese", "Steven Spielberg", "Francis Ford Coppola", "George Lucas"]', 'Steven Spielberg', 1, 'https://example.com/spielberg.jpg', 'Steven Spielberg directed "Jaws", released in 1975, which is considered one of the first summer blockbusters.'),
(6, 'What is the highest-grossing film of all time?', '["Titanic", "Avengers: Endgame", "Avatar", "Star Wars: The Force Awakens"]', 'Avatar', 1, 'https://example.com/avatar.jpg', 'James Cameron''s "Avatar" is the highest-grossing film of all time, having retaken the top spot from "Avengers: Endgame".'),
(6, 'Who played Jack in the movie "Titanic"?', '["Brad Pitt", "Leonardo DiCaprio", "Johnny Depp", "Matt Damon"]', 'Leonardo DiCaprio', 1, 'https://example.com/dicaprio_titanic.jpg', 'Leonardo DiCaprio played Jack Dawson in James Cameron''s "Titanic" (1997).'),
(6, 'Which movie features the character Buzz Lightyear?', '["Shrek", "Toy Story", "Finding Nemo", "The Incredibles"]', 'Toy Story', 1, 'https://example.com/buzz_lightyear.jpg', 'Buzz Lightyear is one of the main characters in the "Toy Story" franchise, first appearing in the 1995 Pixar film.'),
(6, 'Who directed the "Lord of the Rings" trilogy?', '["Christopher Nolan", "James Cameron", "Peter Jackson", "Steven Spielberg"]', 'Peter Jackson', 2, 'https://example.com/peter_jackson.jpg', 'Peter Jackson directed the entire "Lord of the Rings" trilogy, released from 2001 to 2003.'),
(6, 'Which actor has won the most Academy Awards for acting?', '["Meryl Streep", "Jack Nicholson", "Daniel Day-Lewis", "Katharine Hepburn"]', 'Katharine Hepburn', 2, 'https://example.com/katharine_hepburn.jpg', 'Katharine Hepburn won four Academy Awards for Best Actress, more than any other actor or actress.'),
(6, 'What was the first feature-length animated movie?', '["Snow White and the Seven Dwarfs", "Pinocchio", "Fantasia", "Bambi"]', 'Snow White and the Seven Dwarfs', 2, 'https://example.com/snow_white.jpg', 'Disney''s "Snow White and the Seven Dwarfs", released in 1937, was the first full-length cel animated feature film.'),
(6, 'Who played the Joker in "The Dark Knight"?', '["Jack Nicholson", "Heath Ledger", "Jared Leto", "Joaquin Phoenix"]', 'Heath Ledger', 3, 'https://example.com/heath_ledger_joker.jpg', 'Heath Ledger portrayed the Joker in Christopher Nolan''s "The Dark Knight" (2008), winning a posthumous Academy Award for his performance.'),
(6, 'What is the name of the fictional metal in the "X-Men" universe?', '["Vibranium", "Adamantium", "Unobtanium", "Kryptonite"]', 'Adamantium', 3, 'https://example.com/adamantium.jpg', 'Adamantium is the indestructible metal alloy that covers Wolverine''s skeleton and claws in the X-Men series.'),
(6, 'Which film won the first Academy Award for Best Animated Feature?', '["Toy Story", "Shrek", "Monsters, Inc.", "Finding Nemo"]', 'Shrek', 3, 'https://example.com/shrek_oscar.jpg', 'The Academy Award for Best Animated Feature was first awarded in 2002, with "Shrek" winning the inaugural award.');

-- 7. Music
INSERT INTO questions (level_id, question, options, correct_answer, difficulty, image_url, explanation) VALUES
(7, 'Who is known as the "King of Pop"?', '["Elvis Presley", "Michael Jackson", "Prince", "Justin Timberlake"]', 'Michael Jackson', 1, 'https://example.com/michael_jackson.jpg', 'Michael Jackson, known for his significant contributions to music, dance, and fashion, earned the nickname "King of Pop".'),
(7, 'Which band performed the hit song "Bohemian Rhapsody"?', '["The Beatles", "Led Zeppelin", "Queen", "Pink Floyd"]', 'Queen', 1, 'https://example.com/queen.jpg', '"Bohemian Rhapsody" was written by Freddie Mercury for the band Queen''s 1975 album "A Night at the Opera".'),
(7, 'What instrument does a pianist play?', '["Guitar", "Drums", "Piano", "Violin"]', 'Piano', 1, 'https://example.com/piano.jpg', 'A pianist is a musician who plays the piano, a keyboard instrument.'),
(7, 'Who sang "Like a Rolling Stone"?', '["Bob Dylan", "Bruce Springsteen", "Neil Young", "Paul Simon"]', 'Bob Dylan', 1, 'https://example.com/bob_dylan.jpg', 'Bob Dylan wrote and performed "Like a Rolling Stone", released in 1965.'),
(7, 'Which of these is not a wind instrument?', '["Flute", "Clarinet", "Violin", "Saxophone"]', 'Violin', 2, 'https://example.com/violin.jpg', 'The violin is a string instrument, while the others are wind instruments.'),
(7, 'Who was the lead singer of Nirvana?', '["Eddie Vedder", "Kurt Cobain", "Dave Grohl", "Chris Cornell"]', 'Kurt Cobain', 2, 'https://example.com/kurt_cobain.jpg', 'Kurt Cobain was the lead vocalist, guitarist, and primary songwriter of the rock band Nirvana.'),
(7, 'Which composer wrote "The Four Seasons"?', '["Johann Sebastian Bach", "Wolfgang Amadeus Mozart", "Ludwig van Beethoven", "Antonio Vivaldi"]', 'Antonio Vivaldi', 2, 'https://example.com/vivaldi.jpg', 'Antonio Vivaldi composed "The Four Seasons", a group of four violin concerti, around 1720.'),
(7, 'What is the time signature of a standard waltz?', '["2/4", "3/4", "4/4", "6/8"]', '3/4', 3, 'https://example.com/waltz.jpg', 'A waltz is typically in 3/4 time, with three beats per measure and the quarter note getting one beat.'),
(7, 'Which music technology replaced the 8-track tape?', '["Vinyl records", "Cassette tapes", "CDs", "MP3s"]', 'Cassette tapes', 3, 'https://example.com/cassette.jpg', 'Cassette tapes became popular in the 1970s, gradually replacing 8-track tapes before being succeeded by CDs.'),
(7, 'Who composed the opera "The Ring of the Nibelung"?', '["Giuseppe Verdi", "Richard Wagner", "Wolfgang Amadeus Mozart", "Giacomo Puccini"]', 'Richard Wagner', 3, 'https://example.com/wagner.jpg', 'Richard Wagner composed the four-opera cycle "The Ring of the Nibelung" over the course of about 26 years.');

-- 8. Sports
INSERT INTO questions (level_id, question, options, correct_answer, difficulty, image_url, explanation) VALUES
(8, 'In which sport would you perform a slam dunk?', '["Tennis", "Basketball", "Soccer", "Golf"]', 'Basketball', 1, 'https://example.com/slam_dunk.jpg', 'A slam dunk is a type of basketball shot that is performed when a player jumps in the air and manually powers the ball downward through the basket.'),
(8, 'How many players are on a standard soccer team on the field?', '["9", "10", "11", "12"]', '11', 1, 'https://example.com/soccer_team.jpg', 'A standard soccer (football) team has 11 players on the field during a match.'),
(8, 'In which sport is the Stanley Cup awarded?', '["Baseball", "Basketball", "Ice Hockey", "American Football"]', 'Ice Hockey', 1, 'https://example.com/stanley_cup.jpg', 'The Stanley Cup is the championship trophy awarded annually to the playoff champion of the National Hockey League (NHL).'),
(8, 'How many holes are played in a typical round of golf?', '["9", "18", "27", "36"]', '18', 1, 'https://example.com/golf_course.jpg', 'A standard round of golf consists of 18 holes.'),
(8, 'Which country has won the most FIFA World Cup titles?', '["Germany", "Italy", "Argentina", "Brazil"]', 'Brazil', 2, 'https://example.com/brazil_soccer.jpg', 'Brazil has won the FIFA World Cup a record five times (1958, 1962, 1970, 1994, 2002).'),
(8, 'In tennis, what is a "Grand Slam"?', '["Winning all four major tournaments in a calendar year", "Scoring four points in a row", "Winning a match without losing a game", "Beating the top four ranked players"]', 'Winning all four major tournaments in a calendar year', 2, 'https://example.com/grand_slam.jpg', 'A Grand Slam in tennis refers to winning all four major tournaments (Australian Open, French Open, Wimbledon, US Open) in a single calendar year.'),
(8, 'Which swimmer has won the most Olympic gold medals?', '["Ian Thorpe", "Mark Spitz", "Michael Phelps", "Ryan Lochte"]', 'Michael Phelps', 2, 'https://example.com/michael_phelps.jpg', 'Michael Phelps has won a total of 23 Olympic gold medals, the most in Olympic history.'),
(8, 'In cricket, what is a "googly"?', '["A type of bat", "A fielding position", "A ball bowled by a right-arm leg spin bowler that spins from off to leg", "A scoring shot"]', 'A ball bowled by a right-arm leg spin bowler that spins from off to leg', 3, 'https://example.com/googly.jpg', 'A googly is a type of delivery bowled by a right-arm leg spin bowler which spins from the off side to the leg side for a right-handed batsman.'),
(8, 'What is the "fosbury flop"?', '["A swimming stroke", "A gymnastics move", "A high jump technique", "A figure skating jump"]', 'A high jump technique', 3, 'https://example.com/fosbury_flop.jpg', 'The Fosbury Flop is a high jump technique where the athlete goes over the bar backwards, popularized by Dick Fosbury in the 1968 Olympics.'),
(8, 'In baseball, what is a "perfect game"?', '["When a team wins by 10 runs", "When a pitcher strikes out every batter", "When a team scores in every inning", "When no opposing player reaches base"]', 'When no opposing player reaches base', 3, 'https://example.com/perfect_game.jpg', 'A perfect game is when a pitcher (or combination of pitchers) retires every opposing batter without allowing any to reach base over the course of a complete game.');

-- 9. Technology
INSERT INTO questions (level_id, question, options, correct_answer, difficulty, image_url, explanation) VALUES
(9, 'What does "CPU" stand for?', '["Central Processing Unit", "Computer Personal Unit", "Central Processor Unifier", "Computer Processing Utility"]', 'Central Processing Unit', 1, 'https://example.com/cpu.jpg', 'CPU stands for Central Processing Unit, often called the "brain" of the computer.'),
(9, 'Which company created the iPhone?', '["Microsoft", "Google", "Apple", "Samsung"]', 'Apple', 1, 'https://example.com/iphone.jpg', 'The iPhone was created by Apple Inc. and was first released in 2007.'),
(9, 'What does "WWW" stand for in a website browser?', '["World Wide Web", "Western Washington World", "Wide Width Wickets", "World Weather Watch"]', 'World Wide Web', 1, 'https://example.com/www.jpg', 'WWW stands for World Wide Web, the primary system for accessing information over the Internet.'),
(9, 'Which of these is not a programming language?', '["Java", "Python", "Banana", "C++"]', 'Banana', 1, 'https://example.com/programming_languages.jpg', 'Java, Python, and C++ are all programming languages, while Banana is not.'),
(9, 'What is the name of the world's largest social media platform?', '["Twitter", "Instagram", "Facebook", "LinkedIn"]', 'Facebook', 2, 'https://example.com/facebook.jpg', 'Facebook, founded by Mark Zuckerberg, is the world''s largest social media platform with over 2.7 billion monthly active users as of 2021.'),
(9, 'What does "AI" stand for in the tech world?', '["Automated Intelligence", "Artificial Intelligence", "Automated Information", "Artificial Imagery"]', 'Artificial Intelligence', 2, 'https://example.com/ai.jpg', 'AI stands for Artificial Intelligence, which refers to the simulation of human intelligence in machines.'),
(9, 'Which company developed the Android operating system?', '["Apple", "Microsoft", "Google", "Samsung"]', 'Google', 2, 'https://example.com/android.jpg', 'Android was developed by Android Inc., which Google bought in 2005. Google is the primary developer of the Android operating system.'),
(9, 'What is quantum computing?', '["Computing using water", "Computing using light", "Computing using subatomic particles", "Computing using sound waves"]', 'Computing using subatomic particles', 3, 'https://example.com/quantum_computing.jpg', 'Quantum computing uses quantum-mechanical phenomena, such as superposition and entanglement, to perform computation.'),
(9, 'What is the name of the first commercially successful graphical user interface (GUI) operating system?', '["Windows 1.0", "Mac OS", "Linux", "MS-DOS"]', 'Mac OS', 3, 'https://example.com/mac_os.jpg', 'The first commercially successful GUI operating system was Apple''s Mac OS, introduced with the Macintosh in 1984.'),
(9, 'What is the difference between HTTP and HTTPS?', '["Speed", "Security", "Bandwidth", "File size limit"]', 'Security', 3, 'https://example.com/https.jpg', 'HTTPS (Hypertext Transfer Protocol Secure) is the secure version of HTTP, encrypting the data sent between your browser and the website you''re connected to.');

-- 10. Art
INSERT INTO questions (level_id, question, options, correct_answer, difficulty, image_url, explanation) VALUES
(10, 'Who painted the Mona Lisa?', '["Vincent van Gogh", "Leonardo da Vinci", "Pablo Picasso", "Michelangelo"]', 'Leonardo da Vinci', 1, 'https://example.com/mona_lisa.jpg', 'The Mona Lisa was painted by Italian Renaissance artist Leonardo da Vinci in the early 16th century.'),
(10, 'Which of these is a primary color?', '["Green", "Orange", "Purple", "Blue"]', 'Blue', 1, 'https://example.com/primary_colors.jpg', 'Blue is one of the three primary colors, along with red and yellow.'),
(10, 'What type of paint dries the slowest?', '["Acrylic", "Watercolor", "Oil", "Tempera"]', 'Oil', 1, 'https://example.com/oil_paint.jpg', 'Oil paint dries much slower than other types of paint, allowing for longer working times and blending techniques.'),
(10, 'Who cut off his own ear?', '["Claude Monet", "Vincent van Gogh", "Salvador Dali", "Pablo Picasso"]', 'Vincent van Gogh', 1, 'https://example.com/van_gogh.jpg', 'Dutch post-impressionist painter Vincent van Gogh famously cut off a portion of his own ear in 1888.'),
(10, 'Which art movement is Salvador Dali associated with?', '["Impressionism", "Cubism", "Surrealism", "Abstract Expressionism"]', 'Surrealism', 2, 'https://example.com/dali.jpg', 'Salvador Dali was a prominent figure in the Surrealist movement, known for his bizarre and dreamlike images.'),
(10, 'What is chiaroscuro?', '["A type of paint", "A painting technique using strong contrasts between light and dark", "A famous art school", "A sculpture method"]', 'A painting technique using strong contrasts between light and dark', 2, 'https://example.com/chiaroscuro.jpg', 'Chiaroscuro is an art technique that uses strong contrasts between light and dark to create a sense of volume in modeling three-dimensional objects.'),
(10, 'Who painted "The Starry Night"?', '["Claude Monet", "Vincent van Gogh", "Edvard Munch", "Georgia O''Keeffe"]', 'Vincent van Gogh', 2, 'https://example.com/starry_night.jpg', '"The Starry Night" is one of Vincent van Gogh''s most famous works, painted in 1889 while he was in an asylum in France.'),
(10, 'Which ancient civilization is known for building pyramids?', '["Greeks", "Romans", "Egyptians", "Chinese"]', 'Egyptians', 3, 'https://example.com/egyptian_pyramids.jpg', 'The ancient Egyptians are famously known for building pyramids as tombs for their pharaohs, with the Great Pyramid of Giza being one of the Seven Wonders of the Ancient World.'),
(10, 'What is the name of the art movement characterized by abstract patterns and optical illusions?', '["Impressionism", "Cubism", "Op Art", "Fauvism"]', 'Op Art', 3, 'https://example.com/op_art.jpg', 'Op art, short for optical art, is a style of visual art that uses optical illusions to create an impression of movement or hidden images.'),
(10, 'Who sculpted "David"?', '["Leonardo da Vinci", "Michelangelo", "Donatello", "Raphael"]', 'Michelangelo', 3, 'https://example.com/david_sculpture.jpg', 'The statue of David was created by Italian Renaissance artist Michelangelo between 1501 and 1504.');
    `)
    return err
}