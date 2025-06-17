package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

// The Celebrity struct holds detailed information.
type Celebrity struct {
	Name      string
	FamousFor string
	WikiURL   string
}

// The HTML for the greeting page (no changes from last time).
const greetingHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Greeting</title>
     <style>
        body { font-family: sans-serif; text-align: center; margin-top: 50px; }
        h1 { color: #333; }
        p { font-size: 1.2em; line-height: 1.6; }
        ul { list-style-type: none; padding: 0; }
        li { margin: 5px 0; background-color: #f4f4f4; padding: 10px; border-radius: 4px; text-align: left; }
        a { text-decoration: none; font-weight: bold; }
        .back-link { color: #007bff; text-decoration: none; margin-top: 20px; display: inline-block;}
        .future { margin-top: 20px; font-weight: bold; color: #28a745; }
		.celebs { margin-top: 20px; border-top: 1px solid #ccc; padding-top: 15px; display: inline-block; }
    </style>
</head>
<body>
    <h1>{{.Greeting}}</h1>
    
    <p>You are <strong>{{.Age}}</strong> years old.</p>
    <p>Your birthday has fallen on a Sunday <strong>{{.SundayCount}}</strong> times!</p>

    {{if .SundayDates}}
        <p>Those Sundays were:</p>
        <ul>
            {{range .SundayDates}}
                <li>{{.}}</li>
            {{end}}
        </ul>
    {{end}}

    {{if .NextSundayBirthday}}
        <p class="future">Your next birthday on a Sunday will be on {{.NextSundayBirthday}}.</p>
    {{end}}

	{{if .Celebrities}}
		<div class="celebs">
			<p>You share a birthday with these famous people:</p>
			<ul>
				{{range .Celebrities}}
					<li>
						<a href="{{.WikiURL}}" target="_blank">{{.Name}}</a>
						{{if .FamousFor}} - {{.FamousFor}}{{end}}
					</li>
				{{end}}
			</ul>
		</div>
	{{end}}
    
    <a href="/" class="back-link">Go Back</a>
</body>
</html>
`

// PageData struct (no changes).
type PageData struct {
	Greeting           string
	Age                int
	SundayCount        int
	SundayDates        []string
	NextSundayBirthday string
	Celebrities        []Celebrity
}

// All Go logic functions (greetingHandler, indexHandler, main) have no changes.
// They are flexible enough to work with the new data. For brevity, they are placed at the end.

const indexHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Birthday Greeter</title>
    <style>
        body { font-family: sans-serif; text-align: center; margin-top: 50px; }
        form { display: inline-block; }
        input[type="text"], input[type="date"] { width: 200px; padding: 8px; margin-bottom: 10px; }
        input[type="submit"] { padding: 10px 20px; cursor: pointer; }
    </style>
</head>
<body>
    <h1>Enter Your Details</h1>
    <form action="/greet" method="post">
        <label for="name">Name:</label><br>
        <input type="text" id="name" name="name" required><br><br>
        <label for="birthday">Birthday:</label><br>
        <input type="date" id="birthday" name="birthday" required><br><br>
        <input type="submit" value="Submit">
    </form>
</body>
</html>
`

// --- THE PRIMARY UPDATE IS THE DATA IN THIS MAP ---

var celebrityBirthdays = map[string][]Celebrity{
	// January is now fully detailed
	"Jan-01": {
		{Name: "J. D. Salinger", FamousFor: "Author ('The Catcher in the Rye')", WikiURL: "https://en.wikipedia.org/wiki/J._D._Salinger"},
		{Name: "Frank Langella", FamousFor: "Tony Award-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Frank_Langella"},
		{Name: "Grandmaster Flash", FamousFor: "Pioneering hip-hop DJ", WikiURL: "https://en.wikipedia.org/wiki/Grandmaster_Flash"},
	},
	"Jan-02": {
		{Name: "Isaac Asimov", FamousFor: "Prolific science fiction author", WikiURL: "https://en.wikipedia.org/wiki/Isaac_Asimov"},
		{Name: "Cuba Gooding Jr.", FamousFor: "Oscar-winning actor ('Jerry Maguire')", WikiURL: "https://en.wikipedia.org/wiki/Cuba_Gooding_Jr."},
		{Name: "Taye Diggs", FamousFor: "Actor and singer ('Rent')", WikiURL: "https://en.wikipedia.org/wiki/Taye_Diggs"},
	},
	"Jan-03": {
		{Name: "J.R.R. Tolkien", FamousFor: "Author ('The Lord of the Rings')", WikiURL: "https://en.wikipedia.org/wiki/J._R._R._Tolkien"},
		{Name: "Greta Thunberg", FamousFor: "Climate activist", WikiURL: "https://en.wikipedia.org/wiki/Greta_Thunberg"},
		{Name: "Florence Pugh", FamousFor: "Actress ('Midsommar', 'Little Women')", WikiURL: "https://en.wikipedia.org/wiki/Florence_Pugh"},
	},
	"Jan-04": {
		{Name: "Louis Braille", FamousFor: "Inventor of the Braille reading system", WikiURL: "https://en.wikipedia.org/wiki/Louis_Braille"},
		{Name: "Isaac Newton", FamousFor: "Physicist and mathematician (Laws of Motion)", WikiURL: "https://en.wikipedia.org/wiki/Isaac_Newton"},
		{Name: "Michael Stipe", FamousFor: "Lead singer of R.E.M.", WikiURL: "https://en.wikipedia.org/wiki/Michael_Stipe"},
	},
	"Jan-05": {
		{Name: "Hayao Miyazaki", FamousFor: "Co-founder of Studio Ghibli, animator", WikiURL: "https://en.wikipedia.org/wiki/Hayao_Miyazaki"},
		{Name: "Diane Keaton", FamousFor: "Oscar-winning actress ('Annie Hall')", WikiURL: "https://en.wikipedia.org/wiki/Diane_Keaton"},
		{Name: "Bradley Cooper", FamousFor: "Actor and director ('A Star Is Born')", WikiURL: "https://en.wikipedia.org/wiki/Bradley_Cooper"},
	},
	"Jan-06": {
		{Name: "Joan of Arc", FamousFor: "French national heroine and saint", WikiURL: "https://en.wikipedia.org/wiki/Joan_of_Arc"},
		{Name: "Rowan Atkinson", FamousFor: "Actor and comedian ('Mr. Bean')", WikiURL: "https://en.wikipedia.org/wiki/Rowan_Atkinson"},
		{Name: "Eddie Redmayne", FamousFor: "Oscar-winning actor ('The Theory of Everything')", WikiURL: "https://en.wikipedia.org/wiki/Eddie_Redmayne"},
	},
	"Jan-07": {
		{Name: "Nicolas Cage", FamousFor: "Oscar-winning actor ('Leaving Las Vegas')", WikiURL: "https://en.wikipedia.org/wiki/Nicolas_Cage"},
		{Name: "Lewis Hamilton", FamousFor: "Seven-time Formula One World Champion", WikiURL: "https://en.wikipedia.org/wiki/Lewis_Hamilton"},
		{Name: "Jeremy Renner", FamousFor: "Actor ('Hawkeye' in Marvel films)", WikiURL: "https://en.wikipedia.org/wiki/Jeremy_Renner"},
	},
	"Jan-08": {
		{Name: "Elvis Presley", FamousFor: "The 'King of Rock and Roll'", WikiURL: "https://en.wikipedia.org/wiki/Elvis_Presley"},
		{Name: "Stephen Hawking", FamousFor: "Theoretical physicist and cosmologist", WikiURL: "https://en.wikipedia.org/wiki/Stephen_Hawking"},
		{Name: "David Bowie", FamousFor: "Iconic musician and actor", WikiURL: "https://en.wikipedia.org/wiki/David_Bowie"},
	},
	"Jan-09": {
		{Name: "Kate Middleton", FamousFor: "Princess of Wales", WikiURL: "https://en.wikipedia.org/wiki/Catherine,_Princess_of_Wales"},
		{Name: "Jimmy Page", FamousFor: "Guitarist of Led Zeppelin", WikiURL: "https://en.wikipedia.org/wiki/Jimmy_Page"},
		{Name: "Nina Dobrev", FamousFor: "Actress ('The Vampire Diaries')", WikiURL: "https://en.wikipedia.org/wiki/Nina_Dobrev"},
	},
	"Jan-10": {
		{Name: "George Foreman", FamousFor: "Two-time heavyweight boxing champion", WikiURL: "https://en.wikipedia.org/wiki/George_Foreman"},
		{Name: "Rod Stewart", FamousFor: "Rock and pop singer and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Rod_Stewart"},
		{Name: "Jared Kushner", FamousFor: "Businessman and Senior Advisor to President Trump", WikiURL: "https://en.wikipedia.org/wiki/Jared_Kushner"},
	},
	"Jan-11": {
		{Name: "Alexander Hamilton", FamousFor: "Founding Father of the United States", WikiURL: "https://en.wikipedia.org/wiki/Alexander_Hamilton"},
		{Name: "Mary J. Blige", FamousFor: "Grammy-winning singer, 'Queen of Hip-Hop Soul'", WikiURL: "https://en.wikipedia.org/wiki/Mary_J._Blige"},
		{Name: "Amanda Peet", FamousFor: "Actress ('The Whole Nine Yards')", WikiURL: "https://en.wikipedia.org/wiki/Amanda_Peet"},
	},
	"Jan-12": {
		{Name: "Jeff Bezos", FamousFor: "Founder of Amazon", WikiURL: "https://en.wikipedia.org/wiki/Jeff_Bezos"},
		{Name: "Howard Stern", FamousFor: "Radio and television personality", WikiURL: "https://en.wikipedia.org/wiki/Howard_Stern"},
		{Name: "Zayn Malik", FamousFor: "Singer, former member of One Direction", WikiURL: "https://en.wikipedia.org/wiki/Zayn_Malik"},
	},
	"Jan-13": {
		{Name: "Orlando Bloom", FamousFor: "Actor ('Pirates of the Caribbean', 'Lord of the Rings')", WikiURL: "https://en.wikipedia.org/wiki/Orlando_Bloom"},
		{Name: "Julia Louis-Dreyfus", FamousFor: "Emmy-winning actress ('Seinfeld', 'Veep')", WikiURL: "https://en.wikipedia.org/wiki/Julia_Louis-Dreyfus"},
		{Name: "Liam Hemsworth", FamousFor: "Actor ('The Hunger Games')", WikiURL: "https://en.wikipedia.org/wiki/Liam_Hemsworth"},
	},
	"Jan-14": {
		{Name: "LL Cool J", FamousFor: "Rapper and actor", WikiURL: "https://en.wikipedia.org/wiki/LL_Cool_J"},
		{Name: "Dave Grohl", FamousFor: "Frontman of Foo Fighters, drummer for Nirvana", WikiURL: "https://en.wikipedia.org/wiki/Dave_Grohl"},
		{Name: "Jason Bateman", FamousFor: "Actor ('Arrested Development', 'Ozark')", WikiURL: "https://en.wikipedia.org/wiki/Jason_Bateman"},
	},
	"Jan-15": {
		{Name: "Martin Luther King Jr.", FamousFor: "Civil rights leader", WikiURL: "https://en.wikipedia.org/wiki/Martin_Luther_King_Jr."},
		{Name: "Pitbull", FamousFor: "Rapper and singer ('Mr. Worldwide')", WikiURL: "https://en.wikipedia.org/wiki/Pitbull_(rapper)"},
		{Name: "Drew Brees", FamousFor: "Super Bowl-winning NFL Quarterback", WikiURL: "https://en.wikipedia.org/wiki/Drew_Brees"},
	},
	"Jan-16": {
		{Name: "Kate Moss", FamousFor: "Supermodel", WikiURL: "https://en.wikipedia.org/wiki/Kate_Moss"},
		{Name: "Lin-Manuel Miranda", FamousFor: "Creator of the musical 'Hamilton'", WikiURL: "https://en.wikipedia.org/wiki/Lin-Manuel_Miranda"},
		{Name: "Sade", FamousFor: "Singer-songwriter ('Smooth Operator')", WikiURL: "https://en.wikipedia.org/wiki/Sade_(singer)"},
	},
	"Jan-17": {
		{Name: "Muhammad Ali", FamousFor: "Heavyweight boxing champion, 'The Greatest'", WikiURL: "https://en.wikipedia.org/wiki/Muhammad_Ali"},
		{Name: "Michelle Obama", FamousFor: "Former First Lady of the United States", WikiURL: "https://en.wikipedia.org/wiki/Michelle_Obama"},
		{Name: "Jim Carrey", FamousFor: "Actor and comedian", WikiURL: "https://en.wikipedia.org/wiki/Jim_Carrey"},
	},
	"Jan-18": {
		{Name: "A. A. Milne", FamousFor: "Author, creator of 'Winnie-the-Pooh'", WikiURL: "https://en.wikipedia.org/wiki/A._A._Milne"},
		{Name: "Kevin Costner", FamousFor: "Oscar-winning actor and director", WikiURL: "https://en.wikipedia.org/wiki/Kevin_Costner"},
		{Name: "Jason Segel", FamousFor: "Actor and writer ('Forgetting Sarah Marshall')", WikiURL: "https://en.wikipedia.org/wiki/Jason_Segel"},
	},
	"Jan-19": {
		{Name: "Edgar Allan Poe", FamousFor: "Poet and author of macabre stories", WikiURL: "https://en.wikipedia.org/wiki/Edgar_Allan_Poe"},
		{Name: "Dolly Parton", FamousFor: "Country music icon", WikiURL: "https://en.wikipedia.org/wiki/Dolly_Parton"},
		{Name: "Janis Joplin", FamousFor: "Blues and rock singer", WikiURL: "https://en.wikipedia.org/wiki/Janis_Joplin"},
	},
	"Jan-20": {
		{Name: "Buzz Aldrin", FamousFor: "Apollo 11 astronaut, second person on the Moon", WikiURL: "https://en.wikipedia.org/wiki/Buzz_Aldrin"},
		{Name: "David Lynch", FamousFor: "Filmmaker ('Twin Peaks', 'Mulholland Drive')", WikiURL: "https://en.wikipedia.org/wiki/David_Lynch"},
		{Name: "Bill Maher", FamousFor: "Comedian and political commentator", WikiURL: "https://en.wikipedia.org/wiki/Bill_Maher"},
	},
	"Jan-21": {
		{Name: "Stonewall Jackson", FamousFor: "Confederate general during the American Civil War", WikiURL: "https://en.wikipedia.org/wiki/Stonewall_Jackson"},
		{Name: "Geena Davis", FamousFor: "Oscar-winning actress ('Thelma & Louise')", WikiURL: "https://en.wikipedia.org/wiki/Geena_Davis"},
		{Name: "Hakeem Olajuwon", FamousFor: "NBA Hall of Fame basketball player", WikiURL: "https://en.wikipedia.org/wiki/Hakeem_Olajuwon"},
	},
	"Jan-22": {
		{Name: "Lord Byron", FamousFor: "Romantic poet", WikiURL: "https://en.wikipedia.org/wiki/Lord_Byron"},
		{Name: "Diane Lane", FamousFor: "Actress ('Unfaithful')", WikiURL: "https://en.wikipedia.org/wiki/Diane_Lane"},
		{Name: "Guy Fieri", FamousFor: "Chef and television personality ('Diners, Drive-Ins and Dives')", WikiURL: "https://en.wikipedia.org/wiki/Guy_Fieri"},
	},
	"Jan-23": {
		{Name: "John Hancock", FamousFor: "American Revolution leader, famous signature", WikiURL: "https://en.wikipedia.org/wiki/John_Hancock"},
		{Name: "Mariska Hargitay", FamousFor: "Actress ('Law & Order: SVU')", WikiURL: "https://en.wikipedia.org/wiki/Mariska_Hargitay"},
		{Name: "Tito Ortiz", FamousFor: "UFC Hall of Fame mixed martial artist", WikiURL: "https://en.wikipedia.org/wiki/Tito_Ortiz"},
	},
	"Jan-24": {
		{Name: "John Belushi", FamousFor: "Comedian and actor ('Saturday Night Live', 'The Blues Brothers')", WikiURL: "https://en.wikipedia.org/wiki/John_Belushi"},
		{Name: "Neil Diamond", FamousFor: "Singer-songwriter ('Sweet Caroline')", WikiURL: "https://en.wikipedia.org/wiki/Neil_Diamond"},
		{Name: "Mischa Barton", FamousFor: "Actress ('The O.C.')", WikiURL: "https://en.wikipedia.org/wiki/Mischa_Barton"},
	},
	"Jan-25": {
		{Name: "Virginia Woolf", FamousFor: "Modernist author", WikiURL: "https://en.wikipedia.org/wiki/Virginia_Woolf"},
		{Name: "Alicia Keys", FamousFor: "Grammy-winning singer and pianist", WikiURL: "https://en.wikipedia.org/wiki/Alicia_Keys"},
		{Name: "Robert Burns", FamousFor: "National poet of Scotland", WikiURL: "https://en.wikipedia.org/wiki/Robert_Burns"},
	},
	"Jan-26": {
		{Name: "Paul Newman", FamousFor: "Oscar-winning actor and philanthropist", WikiURL: "https://en.wikipedia.org/wiki/Paul_Newman"},
		{Name: "Ellen DeGeneres", FamousFor: "Comedian and television host", WikiURL: "https://en.wikipedia.org/wiki/Ellen_DeGeneres"},
		{Name: "Wayne Gretzky", FamousFor: "Legendary ice hockey player, 'The Great One'", WikiURL: "https://en.wikipedia.org/wiki/Wayne_Gretzky"},
	},
	"Jan-27": {
		{Name: "Wolfgang Amadeus Mozart", FamousFor: "Classical composer", WikiURL: "https://en.wikipedia.org/wiki/Wolfgang_Amadeus_Mozart"},
		{Name: "Lewis Carroll", FamousFor: "Author ('Alice's Adventures in Wonderland')", WikiURL: "https://en.wikipedia.org/wiki/Lewis_Carroll"},
		{Name: "Chris Rock", FamousFor: "Comedian and actor", WikiURL: "https://en.wikipedia.org/wiki/Chris_Rock"},
	},
	"Jan-28": {
		{Name: "Elijah Wood", FamousFor: "Actor ('Frodo' in 'The Lord of the Rings')", WikiURL: "https://en.wikipedia.org/wiki/Elijah_Wood"},
		{Name: "J. Cole", FamousFor: "Grammy-winning rapper and producer", WikiURL: "https://en.wikipedia.org/wiki/J._Cole"},
		{Name: "Nicolas Sarkozy", FamousFor: "Former President of France", WikiURL: "https://en.wikipedia.org/wiki/Nicolas_Sarkozy"},
	},
	"Jan-29": {
		{Name: "Oprah Winfrey", FamousFor: "Media executive and talk show host", WikiURL: "https://en.wikipedia.org/wiki/Oprah_Winfrey"},
		{Name: "Tom Selleck", FamousFor: "Actor ('Magnum, P.I.')", WikiURL: "https://en.wikipedia.org/wiki/Tom_Selleck"},
		{Name: "Adam Lambert", FamousFor: "Singer, frontman for Queen + Adam Lambert", WikiURL: "https://en.wikipedia.org/wiki/Adam_Lambert"},
	},
	"Jan-30": {
		{Name: "Franklin D. Roosevelt", FamousFor: "32nd U.S. President", WikiURL: "https://en.wikipedia.org/wiki/Franklin_D._Roosevelt"},
		{Name: "Christian Bale", FamousFor: "Oscar-winning actor ('Batman')", WikiURL: "https://en.wikipedia.org/wiki/Christian_Bale"},
		{Name: "Gene Hackman", FamousFor: "Oscar-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Gene_Hackman"},
	},
	"Jan-31": {
		{Name: "Jackie Robinson", FamousFor: "First African American to play in Major League Baseball", WikiURL: "https://en.wikipedia.org/wiki/Jackie_Robinson"},
		{Name: "Justin Timberlake", FamousFor: "Grammy-winning singer and actor", WikiURL: "https://en.wikipedia.org/wiki/Justin_Timberlake"},
		{Name: "Kerry Washington", FamousFor: "Actress ('Scandal')", WikiURL: "https://en.wikipedia.org/wiki/Kerry_Washington"},
	},
	"Feb-01": {
		{Name: "Harry Styles", FamousFor: "Grammy-winning singer, former member of One Direction", WikiURL: "https://en.wikipedia.org/wiki/Harry_Styles"},
		{Name: "Clark Gable", FamousFor: "Actor, 'The King of Hollywood' ('Gone with the Wind')", WikiURL: "https://en.wikipedia.org/wiki/Clark_Gable"},
		{Name: "Lisa Marie Presley", FamousFor: "Singer-songwriter, daughter of Elvis Presley", WikiURL: "https://en.wikipedia.org/wiki/Lisa_Marie_Presley"},
	},
	"Feb-02": {
		{Name: "Shakira", FamousFor: "Grammy-winning singer ('Hips Don't Lie')", WikiURL: "https://en.wikipedia.org/wiki/Shakira"},
		{Name: "Gerard Piqué", FamousFor: "Professional footballer (FC Barcelona)", WikiURL: "https://en.wikipedia.org/wiki/Gerard_Piqu%C3%A9"},
		{Name: "James Joyce", FamousFor: "Author ('Ulysses')", WikiURL: "https://en.wikipedia.org/wiki/James_Joyce"},
	},
	"Feb-03": {
		{Name: "Norman Rockwell", FamousFor: "Painter and illustrator of American culture", WikiURL: "https://en.wikipedia.org/wiki/Norman_Rockwell"},
		{Name: "Amal Clooney", FamousFor: "Human rights lawyer", WikiURL: "https://en.wikipedia.org/wiki/Amal_Clooney"},
		{Name: "Nathan Lane", FamousFor: "Tony-winning actor ('The Producers', 'The Lion King')", WikiURL: "https://en.wikipedia.org/wiki/Nathan_Lane"},
	},
	"Feb-04": {
		{Name: "Rosa Parks", FamousFor: "Civil rights activist", WikiURL: "https://en.wikipedia.org/wiki/Rosa_Parks"},
		{Name: "Alice Cooper", FamousFor: "Rock singer and musician", WikiURL: "https://en.wikipedia.org/wiki/Alice_Cooper"},
		{Name: "Oscar De La Hoya", FamousFor: "Olympic gold medalist professional boxer", WikiURL: "https://en.wikipedia.org/wiki/Oscar_De_La_Hoya"},
	},
	"Feb-05": {
		{Name: "Cristiano Ronaldo", FamousFor: "Professional footballer, one of the greatest of all time", WikiURL: "https://en.wikipedia.org/wiki/Cristiano_Ronaldo"},
		{Name: "Neymar", FamousFor: "Brazilian professional footballer", WikiURL: "https://en.wikipedia.org/wiki/Neymar"},
		{Name: "Hank Aaron", FamousFor: "Hall of Fame baseball player", WikiURL: "https://en.wikipedia.org/wiki/Hank_Aaron"},
	},
	"Feb-06": {
		{Name: "Bob Marley", FamousFor: "Reggae music icon", WikiURL: "https://en.wikipedia.org/wiki/Bob_Marley"},
		{Name: "Ronald Reagan", FamousFor: "40th U.S. President", WikiURL: "https://en.wikipedia.org/wiki/Ronald_Reagan"},
		{Name: "Axl Rose", FamousFor: "Lead vocalist of Guns N' Roses", WikiURL: "https://en.wikipedia.org/wiki/Axl_Rose"},
	},
	"Feb-07": {
		{Name: "Charles Dickens", FamousFor: "Author ('A Christmas Carol', 'Oliver Twist')", WikiURL: "https://en.wikipedia.org/wiki/Charles_Dickens"},
		{Name: "Chris Rock", FamousFor: "Comedian and actor", WikiURL: "https://en.wikipedia.org/wiki/Chris_Rock"},
		{Name: "Ashton Kutcher", FamousFor: "Actor and investor", WikiURL: "https://en.wikipedia.org/wiki/Ashton_Kutcher"},
	},
	"Feb-08": {
		{Name: "Jules Verne", FamousFor: "Author ('Twenty Thousand Leagues Under the Seas')", WikiURL: "https://en.wikipedia.org/wiki/Jules_Verne"},
		{Name: "John Williams", FamousFor: "Composer of film scores ('Star Wars', 'Jurassic Park')", WikiURL: "https://en.wikipedia.org/wiki/John_Williams"},
		{Name: "James Dean", FamousFor: "Actor and cultural icon ('Rebel Without a Cause')", WikiURL: "https://en.wikipedia.org/wiki/James_Dean"},
	},
	"Feb-09": {
		{Name: "Joe Pesci", FamousFor: "Oscar-winning actor ('Goodfellas')", WikiURL: "https://en.wikipedia.org/wiki/Joe_Pesci"},
		{Name: "Tom Hiddleston", FamousFor: "Actor ('Loki' in Marvel films)", WikiURL: "https://en.wikipedia.org/wiki/Tom_Hiddleston"},
		{Name: "Michael B. Jordan", FamousFor: "Actor ('Black Panther', 'Creed')", WikiURL: "https://en.wikipedia.org/wiki/Michael_B._Jordan"},
	},
	"Feb-10": {
		{Name: "Mark Spitz", FamousFor: "Olympic swimmer, 9-time gold medalist", WikiURL: "https://en.wikipedia.org/wiki/Mark_Spitz"},
		{Name: "Laura Dern", FamousFor: "Oscar-winning actress ('Jurassic Park')", WikiURL: "https://en.wikipedia.org/wiki/Laura_Dern"},
		{Name: "Elizabeth Banks", FamousFor: "Actress and director ('The Hunger Games')", WikiURL: "https://en.wikipedia.org/wiki/Elizabeth_Banks"},
	},
	"Feb-11": {
		{Name: "Thomas Edison", FamousFor: "Inventor (light bulb, phonograph)", WikiURL: "https://en.wikipedia.org/wiki/Thomas_Edison"},
		{Name: "Jennifer Aniston", FamousFor: "Actress ('Rachel' in 'Friends')", WikiURL: "https://en.wikipedia.org/wiki/Jennifer_Aniston"},
		{Name: "Taylor Lautner", FamousFor: "Actor ('Twilight' series)", WikiURL: "https://en.wikipedia.org/wiki/Taylor_Lautner"},
	},
	"Feb-12": {
		{Name: "Abraham Lincoln", FamousFor: "16th U.S. President", WikiURL: "https://en.wikipedia.org/wiki/Abraham_Lincoln"},
		{Name: "Charles Darwin", FamousFor: "Naturalist, theory of evolution", WikiURL: "https://en.wikipedia.org/wiki/Charles_Darwin"},
		{Name: "Christina Ricci", FamousFor: "Actress ('The Addams Family', 'Wednesday')", WikiURL: "https://en.wikipedia.org/wiki/Christina_Ricci"},
	},
	"Feb-13": {
		{Name: "Peter Gabriel", FamousFor: "Musician, original lead singer of Genesis", WikiURL: "https://en.wikipedia.org/wiki/Peter_Gabriel"},
		{Name: "Robbie Williams", FamousFor: "Singer-songwriter", WikiURL: "https://en.wikipedia.org/wiki/Robbie_Williams"},
		{Name: "Jerry Springer", FamousFor: "Television host ('The Jerry Springer Show')", WikiURL: "https://en.wikipedia.org/wiki/Jerry_Springer"},
	},
	"Feb-14": {
		{Name: "Michael Bloomberg", FamousFor: "Businessman, former Mayor of New York City", WikiURL: "https://en.wikipedia.org/wiki/Michael_Bloomberg"},
		{Name: "Simon Pegg", FamousFor: "Actor and writer ('Shaun of the Dead', 'Hot Fuzz')", WikiURL: "https://en.wikipedia.org/wiki/Simon_Pegg"},
		{Name: "Frederick Douglass", FamousFor: "Abolitionist and writer", WikiURL: "https://en.wikipedia.org/wiki/Frederick_Douglass"},
	},
	"Feb-15": {
		{Name: "Galileo Galilei", FamousFor: "Astronomer and physicist", WikiURL: "https://en.wikipedia.org/wiki/Galileo_Galilei"},
		{Name: "Matt Groening", FamousFor: "Creator of 'The Simpsons' and 'Futurama'", WikiURL: "https://en.wikipedia.org/wiki/Matt_Groening"},
		{Name: "Jane Seymour", FamousFor: "Actress ('Dr. Quinn, Medicine Woman')", WikiURL: "https://en.wikipedia.org/wiki/Jane_Seymour_(actress)"},
	},
	"Feb-16": {
		{Name: "John McEnroe", FamousFor: "Hall of Fame tennis player", WikiURL: "https://en.wikipedia.org/wiki/John_McEnroe"},
		{Name: "Ice-T", FamousFor: "Rapper and actor ('Law & Order: SVU')", WikiURL: "https://en.wikipedia.org/wiki/Ice-T"},
		{Name: "The Weeknd", FamousFor: "Grammy-winning singer and songwriter", WikiURL: "https://en.wikipedia.org/wiki/The_Weeknd"},
	},
	"Feb-17": {
		{Name: "Michael Jordan", FamousFor: "Legendary professional basketball player", WikiURL: "https://en.wikipedia.org/wiki/Michael_Jordan"},
		{Name: "Ed Sheeran", FamousFor: "Grammy-winning singer-songwriter", WikiURL: "https://en.wikipedia.org/wiki/Ed_Sheeran"},
		{Name: "Paris Hilton", FamousFor: "Media personality and socialite", WikiURL: "https://en.wikipedia.org/wiki/Paris_Hilton"},
	},
	"Feb-18": {
		{Name: "John Travolta", FamousFor: "Actor ('Grease', 'Pulp Fiction')", WikiURL: "https://en.wikipedia.org/wiki/John_Travolta"},
		{Name: "Dr. Dre", FamousFor: "Hip-hop producer, rapper, and entrepreneur", WikiURL: "https://en.wikipedia.org/wiki/Dr._Dre"},
		{Name: "Yoko Ono", FamousFor: "Artist, musician, wife of John Lennon", WikiURL: "https://en.wikipedia.org/wiki/Yoko_Ono"},
	},
	"Feb-19": {
		{Name: "Nicolaus Copernicus", FamousFor: "Astronomer, formulated a model of the universe with the Sun at the center", WikiURL: "https://en.wikipedia.org/wiki/Nicolaus_Copernicus"},
		{Name: "Smokey Robinson", FamousFor: "Singer, songwriter, 'King of Motown'", WikiURL: "https://en.wikipedia.org/wiki/Smokey_Robinson"},
		{Name: "Millie Bobby Brown", FamousFor: "Actress ('Eleven' in 'Stranger Things')", WikiURL: "https://en.wikipedia.org/wiki/Millie_Bobby_Brown"},
	},
	"Feb-20": {
		{Name: "Rihanna", FamousFor: "Grammy-winning singer and businesswoman", WikiURL: "https://en.wikipedia.org/wiki/Rihanna"},
		{Name: "Kurt Cobain", FamousFor: "Lead singer and guitarist of Nirvana", WikiURL: "https://en.wikipedia.org/wiki/Kurt_Cobain"},
		{Name: "Sidney Poitier", FamousFor: "First African American to win Oscar for Best Actor", WikiURL: "https://en.wikipedia.org/wiki/Sidney_Poitier"},
	},
	"Feb-21": {
		{Name: "Nina Simone", FamousFor: "Singer, songwriter, and civil rights activist", WikiURL: "https://en.wikipedia.org/wiki/Nina_Simone"},
		{Name: "Alan Rickman", FamousFor: "Actor ('Snape' in 'Harry Potter')", WikiURL: "https://en.wikipedia.org/wiki/Alan_Rickman"},
		{Name: "Elliot Page", FamousFor: "Actor ('Juno', 'The Umbrella Academy')", WikiURL: "https://en.wikipedia.org/wiki/Elliot_Page"},
	},
	"Feb-22": {
		{Name: "George Washington", FamousFor: "First U.S. President", WikiURL: "https://en.wikipedia.org/wiki/George_Washington"},
		{Name: "Steve Irwin", FamousFor: "Wildlife expert, 'The Crocodile Hunter'", WikiURL: "https://en.wikipedia.org/wiki/Steve_Irwin"},
		{Name: "Drew Barrymore", FamousFor: "Actress and talk show host", WikiURL: "https://en.wikipedia.org/wiki/Drew_Barrymore"},
	},
	"Feb-23": {
		{Name: "W. E. B. Du Bois", FamousFor: "Sociologist, historian, and civil rights activist", WikiURL: "https://en.wikipedia.org/wiki/W._E._B._Du_Bois"},
		{Name: "Dakota Fanning", FamousFor: "Actress ('I Am Sam', 'War of the Worlds')", WikiURL: "https://en.wikipedia.org/wiki/Dakota_Fanning"},
		{Name: "Emily Blunt", FamousFor: "Actress ('The Devil Wears Prada', 'A Quiet Place')", WikiURL: "https://en.wikipedia.org/wiki/Emily_Blunt"},
	},
	"Feb-24": {
		{Name: "Steve Jobs", FamousFor: "Co-founder of Apple Inc.", WikiURL: "https://en.wikipedia.org/wiki/Steve_Jobs"},
		{Name: "Floyd Mayweather Jr.", FamousFor: "Undefeated professional boxer", WikiURL: "https://en.wikipedia.org/wiki/Floyd_Mayweather_Jr."},
		{Name: "Daniel Kaluuya", FamousFor: "Oscar-winning actor ('Get Out', 'Judas and the Black Messiah')", WikiURL: "https://en.wikipedia.org/wiki/Daniel_Kaluuya"},
	},
	"Feb-25": {
		{Name: "George Harrison", FamousFor: "Lead guitarist of The Beatles", WikiURL: "https://en.wikipedia.org/wiki/George_Harrison"},
		{Name: "Sean Astin", FamousFor: "Actor ('Samwise' in 'The Lord of the Rings', 'The Goonies')", WikiURL: "https://en.wikipedia.org/wiki/Sean_Astin"},
		{Name: "Rashida Jones", FamousFor: "Actress ('Parks and Recreation', 'The Office')", WikiURL: "https://en.wikipedia.org/wiki/Rashida_Jones"},
	},
	"Feb-26": {
		{Name: "Johnny Cash", FamousFor: "Country music singer-songwriter, 'The Man in Black'", WikiURL: "https://en.wikipedia.org/wiki/Johnny_Cash"},
		{Name: "Victor Hugo", FamousFor: "Author ('Les Misérables', 'The Hunchback of Notre-Dame')", WikiURL: "https://en.wikipedia.org/wiki/Victor_Hugo"},
		{Name: "Erykah Badu", FamousFor: "Grammy-winning singer, 'Queen of Neo Soul'", WikiURL: "https://en.wikipedia.org/wiki/Erykah_Badu"},
	},
	"Feb-27": {
		{Name: "John Steinbeck", FamousFor: "Author ('The Grapes of Wrath', 'Of Mice and Men')", WikiURL: "https://en.wikipedia.org/wiki/John_Steinbeck"},
		{Name: "Elizabeth Taylor", FamousFor: "Oscar-winning actress and icon", WikiURL: "https://en.wikipedia.org/wiki/Elizabeth_Taylor"},
		{Name: "Kate Mara", FamousFor: "Actress ('House of Cards')", WikiURL: "https://en.wikipedia.org/wiki/Kate_Mara"},
	},
	"Feb-28": {
		{Name: "Mario Andretti", FamousFor: "Formula One and Indianapolis 500 racing champion", WikiURL: "https://en.wikipedia.org/wiki/Mario_Andretti"},
		{Name: "Frank Gehry", FamousFor: "Pritzker Prize-winning architect", WikiURL: "https://en.wikipedia.org/wiki/Frank_Gehry"},
		{Name: "Bernadette Peters", FamousFor: "Tony-winning Broadway actress and singer", WikiURL: "https://en.wikipedia.org/wiki/Bernadette_Peters"},
	},
	"Feb-29": {
		{Name: "Ja Rule", FamousFor: "Rapper and singer", WikiURL: "https://en.wikipedia.org/wiki/Ja_Rule"},
		{Name: "Tony Robbins", FamousFor: "Author and motivational speaker", WikiURL: "https://en.wikipedia.org/wiki/Tony_Robbins"},
		{Name: "Superman (fictional)", FamousFor: "Superhero from Krypton", WikiURL: "https://en.wikipedia.org/wiki/Superman"},
	},
	"Mar-01": {
		{Name: "Frédéric Chopin", FamousFor: "Composer and virtuoso pianist", WikiURL: "https://en.wikipedia.org/wiki/Fr%C3%A9d%C3%A9ric_Chopin"},
		{Name: "Justin Bieber", FamousFor: "Grammy-winning pop singer", WikiURL: "https://en.wikipedia.org/wiki/Justin_Bieber"},
		{Name: "Lupita Nyong'o", FamousFor: "Oscar-winning actress ('12 Years a Slave')", WikiURL: "https://en.wikipedia.org/wiki/Lupita_Nyong%27o"},
	},
	"Mar-02": {
		{Name: "Dr. Seuss", FamousFor: "Children's author ('The Cat in the Hat')", WikiURL: "https://en.wikipedia.org/wiki/Dr._Seuss"},
		{Name: "Jon Bon Jovi", FamousFor: "Rock singer, frontman of Bon Jovi", WikiURL: "https://en.wikipedia.org/wiki/Jon_Bon_Jovi"},
		{Name: "Daniel Craig", FamousFor: "Actor ('James Bond')", WikiURL: "https://en.wikipedia.org/wiki/Daniel_Craig"},
	},
	"Mar-03": {
		{Name: "Alexander Graham Bell", FamousFor: "Inventor of the telephone", WikiURL: "https://en.wikipedia.org/wiki/Alexander_Graham_Bell"},
		{Name: "Jessica Biel", FamousFor: "Actress", WikiURL: "https://en.wikipedia.org/wiki/Jessica_Biel"},
		{Name: "Camila Cabello", FamousFor: "Singer and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Camila_Cabello"},
	},
	"Mar-04": {
		{Name: "Antonio Vivaldi", FamousFor: "Composer ('The Four Seasons')", WikiURL: "https://en.wikipedia.org/wiki/Antonio_Vivaldi"},
		{Name: "Catherine O'Hara", FamousFor: "Emmy-winning actress ('Schitt's Creek')", WikiURL: "https://en.wikipedia.org/wiki/Catherine_O%27Hara"},
		{Name: "Patricia Heaton", FamousFor: "Actress ('Everybody Loves Raymond')", WikiURL: "https://en.wikipedia.org/wiki/Patricia_Heaton"},
	},
	"Mar-05": {
		{Name: "Eva Mendes", FamousFor: "Actress and model", WikiURL: "https://en.wikipedia.org/wiki/Eva_Mendes"},
		{Name: "Penn Jillette", FamousFor: "Magician, part of Penn & Teller", WikiURL: "https://en.wikipedia.org/wiki/Penn_Jillette"},
		{Name: "Madison Beer", FamousFor: "Singer and media personality", WikiURL: "https://en.wikipedia.org/wiki/Madison_Beer"},
	},
	"Mar-06": {
		{Name: "Michelangelo", FamousFor: "Renaissance artist (Sistine Chapel)", WikiURL: "https://en.wikipedia.org/wiki/Michelangelo"},
		{Name: "Shaquille O'Neal", FamousFor: "Hall of Fame basketball player", WikiURL: "https://en.wikipedia.org/wiki/Shaquille_O%27Neal"},
		{Name: "Rob Reiner", FamousFor: "Director and actor ('The Princess Bride')", WikiURL: "https://en.wikipedia.org/wiki/Rob_Reiner"},
	},
	"Mar-07": {
		{Name: "Bryan Cranston", FamousFor: "Emmy-winning actor ('Walter White' in 'Breaking Bad')", WikiURL: "https://en.wikipedia.org/wiki/Bryan_Cranston"},
		{Name: "Rachel Weisz", FamousFor: "Oscar-winning actress ('The Mummy')", WikiURL: "https://en.wikipedia.org/wiki/Rachel_Weisz"},
		{Name: "Jenna Fischer", FamousFor: "Actress ('Pam' in 'The Office')", WikiURL: "https://en.wikipedia.org/wiki/Jenna_Fischer"},
	},
	"Mar-08": {
		{Name: "James Van Der Beek", FamousFor: "Actor ('Dawson's Creek')", WikiURL: "https://en.wikipedia.org/wiki/James_Van_Der_Beek"},
		{Name: "Aidan Quinn", FamousFor: "Actor ('Legends of the Fall')", WikiURL: "https://en.wikipedia.org/wiki/Aidan_Quinn"},
		{Name: "Freddie Prinze Jr.", FamousFor: "Actor ('I Know What You Did Last Summer')", WikiURL: "https://en.wikipedia.org/wiki/Freddie_Prinze_Jr."},
	},
	"Mar-09": {
		{Name: "Amerigo Vespucci", FamousFor: "Explorer, America was named after him", WikiURL: "https://en.wikipedia.org/wiki/Amerigo_Vespucci"},
		{Name: "Juliette Binoche", FamousFor: "Oscar-winning actress ('The English Patient')", WikiURL: "https://en.wikipedia.org/wiki/Juliette_Binoche"},
		{Name: "Oscar Isaac", FamousFor: "Actor ('Star Wars', 'Dune')", WikiURL: "https://en.wikipedia.org/wiki/Oscar_Isaac"},
	},
	"Mar-10": {
		{Name: "Chuck Norris", FamousFor: "Martial artist, actor, and internet meme", WikiURL: "https://en.wikipedia.org/wiki/Chuck_Norris"},
		{Name: "Sharon Stone", FamousFor: "Actress ('Basic Instinct')", WikiURL: "https://en.wikipedia.org/wiki/Sharon_Stone"},
		{Name: "Olivia Wilde", FamousFor: "Actress and director", WikiURL: "https://en.wikipedia.org/wiki/Olivia_Wilde"},
	},
	"Mar-11": {
		{Name: "Rupert Murdoch", FamousFor: "Media mogul (Fox News, The Wall Street Journal)", WikiURL: "https://en.wikipedia.org/wiki/Rupert_Murdoch"},
		{Name: "Johnny Knoxville", FamousFor: "Co-creator and star of 'Jackass'", WikiURL: "https://en.wikipedia.org/wiki/Johnny_Knoxville"},
		{Name: "Thora Birch", FamousFor: "Actress ('American Beauty')", WikiURL: "https://en.wikipedia.org/wiki/Thora_Birch"},
	},
	"Mar-12": {
		{Name: "Liza Minnelli", FamousFor: "Oscar-winning actress and singer", WikiURL: "https://en.wikipedia.org/wiki/Liza_Minnelli"},
		{Name: "James Taylor", FamousFor: "Grammy-winning singer-songwriter", WikiURL: "https://en.wikipedia.org/wiki/James_Taylor"},
		{Name: "Mitt Romney", FamousFor: "U.S. Senator and former presidential candidate", WikiURL: "https://en.wikipedia.org/wiki/Mitt_Romney"},
	},
	"Mar-13": {
		{Name: "Neil Sedaka", FamousFor: "Singer and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Neil_Sedaka"},
		{Name: "Common", FamousFor: "Grammy and Oscar-winning rapper and actor", WikiURL: "https://en.wikipedia.org/wiki/Common_(rapper)"},
		{Name: "Kaya Scodelario", FamousFor: "Actress ('Skins', 'The Maze Runner')", WikiURL: "https://en.wikipedia.org/wiki/Kaya_Scodelario"},
	},
	"Mar-14": {
		{Name: "Albert Einstein", FamousFor: "Theoretical physicist (Theory of Relativity)", WikiURL: "https://en.wikipedia.org/wiki/Albert_Einstein"},
		{Name: "Michael Caine", FamousFor: "Oscar-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Michael_Caine"},
		{Name: "Stephen Curry", FamousFor: "Championship-winning basketball player", WikiURL: "https://en.wikipedia.org/wiki/Stephen_Curry"},
	},
	"Mar-15": {
		{Name: "Julius Caesar", FamousFor: "Roman general and statesman", WikiURL: "https://en.wikipedia.org/wiki/Julius_Caesar"},
		{Name: "will.i.am", FamousFor: "Musician, frontman of The Black Eyed Peas", WikiURL: "https://en.wikipedia.org/wiki/Will.i.am"},
		{Name: "Eva Longoria", FamousFor: "Actress ('Desperate Housewives')", WikiURL: "https://en.wikipedia.org/wiki/Eva_Longoria"},
	},
	"Mar-16": {
		{Name: "Jerry Lewis", FamousFor: "Comedian and actor", WikiURL: "https://en.wikipedia.org/wiki/Jerry_Lewis"},
		{Name: "Blake Griffin", FamousFor: "Professional basketball player", WikiURL: "https://en.wikipedia.org/wiki/Blake_Griffin"},
		{Name: "Lauren Graham", FamousFor: "Actress ('Gilmore Girls')", WikiURL: "https://en.wikipedia.org/wiki/Lauren_Graham"},
	},
	"Mar-17": {
		{Name: "St. Patrick", FamousFor: "Patron saint of Ireland", WikiURL: "https://en.wikipedia.org/wiki/Saint_Patrick"},
		{Name: "Kurt Russell", FamousFor: "Actor ('Escape from New York')", WikiURL: "https://en.wikipedia.org/wiki/Kurt_Russell"},
		{Name: "Rob Lowe", FamousFor: "Actor ('The West Wing', 'Parks and Recreation')", WikiURL: "https://en.wikipedia.org/wiki/Rob_Lowe"},
	},
	"Mar-18": {
		{Name: "Adam Levine", FamousFor: "Lead singer of Maroon 5", WikiURL: "https://en.wikipedia.org/wiki/Adam_Levine"},
		{Name: "Queen Latifah", FamousFor: "Rapper, singer, and actress", WikiURL: "https://en.wikipedia.org/wiki/Queen_Latifah"},
		{Name: "Lily Collins", FamousFor: "Actress ('Emily in Paris')", WikiURL: "https://en.wikipedia.org/wiki/Lily_Collins"},
	},
	"Mar-19": {
		{Name: "Bruce Willis", FamousFor: "Actor ('Die Hard' series)", WikiURL: "https://en.wikipedia.org/wiki/Bruce_Willis"},
		{Name: "Glenn Close", FamousFor: "Emmy and Tony-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Glenn_Close"},
		{Name: "Harvey Weinstein", FamousFor: "Disgraced film producer", WikiURL: "https://en.wikipedia.org/wiki/Harvey_Weinstein"},
	},
	"Mar-20": {
		{Name: "Mr. Rogers (Fred Rogers)", FamousFor: "Host of 'Mister Rogers' Neighborhood'", WikiURL: "https://en.wikipedia.org/wiki/Fred_Rogers"},
		{Name: "Spike Lee", FamousFor: "Oscar-winning director and producer", WikiURL: "https://en.wikipedia.org/wiki/Spike_Lee"},
		{Name: "Chester Bennington", FamousFor: "Lead vocalist of Linkin Park", WikiURL: "https://en.wikipedia.org/wiki/Chester_Bennington"},
	},
	"Mar-21": {
		{Name: "Johann Sebastian Bach", FamousFor: "Baroque composer", WikiURL: "https://en.wikipedia.org/wiki/Johann_Sebastian_Bach"},
		{Name: "Gary Oldman", FamousFor: "Oscar-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Gary_Oldman"},
		{Name: "Matthew Broderick", FamousFor: "Actor ('Ferris Bueller's Day Off')", WikiURL: "https://en.wikipedia.org/wiki/Matthew_Broderick"},
	},
	"Mar-22": {
		{Name: "William Shatner", FamousFor: "Actor ('Captain Kirk' in 'Star Trek')", WikiURL: "https://en.wikipedia.org/wiki/William_Shatner"},
		{Name: "Reese Witherspoon", FamousFor: "Oscar-winning actress ('Legally Blonde')", WikiURL: "https://en.wikipedia.org/wiki/Reese_Witherspoon"},
		{Name: "Andrew Lloyd Webber", FamousFor: "Composer of musicals ('The Phantom of the Opera')", WikiURL: "https://en.wikipedia.org/wiki/Andrew_Lloyd_Webber"},
	},
	"Mar-23": {
		{Name: "Akira Kurosawa", FamousFor: "Japanese film director ('Seven Samurai')", WikiURL: "https://en.wikipedia.org/wiki/Akira_Kurosawa"},
		{Name: "Chaka Khan", FamousFor: "Grammy-winning singer, 'Queen of Funk'", WikiURL: "https://en.wikipedia.org/wiki/Chaka_Khan"},
		{Name: "Perez Hilton", FamousFor: "Blogger and media personality", WikiURL: "https://en.wikipedia.org/wiki/Perez_Hilton"},
	},
	"Mar-24": {
		{Name: "Harry Houdini", FamousFor: "Escape artist and magician", WikiURL: "https://en.wikipedia.org/wiki/Harry_Houdini"},
		{Name: "Steve McQueen", FamousFor: "Actor, 'The King of Cool'", WikiURL: "https://en.wikipedia.org/wiki/Steve_McQueen"},
		{Name: "Jessica Chastain", FamousFor: "Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Jessica_Chastain"},
	},
	"Mar-25": {
		{Name: "Elton John", FamousFor: "Grammy-winning singer and pianist", WikiURL: "https://en.wikipedia.org/wiki/Elton_John"},
		{Name: "Aretha Franklin", FamousFor: "The 'Queen of Soul'", WikiURL: "https://en.wikipedia.org/wiki/Aretha_Franklin"},
		{Name: "Sarah Jessica Parker", FamousFor: "Actress ('Carrie Bradshaw' in 'Sex and the City')", WikiURL: "https://en.wikipedia.org/wiki/Sarah_Jessica_Parker"},
	},
	"Mar-26": {
		{Name: "Robert Frost", FamousFor: "Poet", WikiURL: "https://en.wikipedia.org/wiki/Robert_Frost"},
		{Name: "Steven Tyler", FamousFor: "Lead singer of Aerosmith", WikiURL: "https://en.wikipedia.org/wiki/Steven_Tyler"},
		{Name: "Keira Knightley", FamousFor: "Actress ('Pirates of the Caribbean')", WikiURL: "https://en.wikipedia.org/wiki/Keira_Knightley"},
	},
	"Mar-27": {
		{Name: "Quentin Tarantino", FamousFor: "Oscar-winning director ('Pulp Fiction')", WikiURL: "https://en.wikipedia.org/wiki/Quentin_Tarantino"},
		{Name: "Mariah Carey", FamousFor: "Grammy-winning singer with a five-octave range", WikiURL: "https://en.wikipedia.org/wiki/Mariah_Carey"},
		{Name: "Jessie J", FamousFor: "Singer and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Jessie_J"},
	},
	"Mar-28": {
		{Name: "Lady Gaga", FamousFor: "Grammy and Oscar-winning singer and actress", WikiURL: "https://en.wikipedia.org/wiki/Lady_Gaga"},
		{Name: "Vince Vaughn", FamousFor: "Actor and comedian ('Wedding Crashers')", WikiURL: "https://en.wikipedia.org/wiki/Vince_Vaughn"},
		{Name: "Reba McEntire", FamousFor: "Country music singer and actress", WikiURL: "https://en.wikipedia.org/wiki/Reba_McEntire"},
	},
	"Mar-29": {
		{Name: "John Tyler", FamousFor: "10th U.S. President", WikiURL: "https://en.wikipedia.org/wiki/John_Tyler"},
		{Name: "Lucy Lawless", FamousFor: "Actress ('Xena: Warrior Princess')", WikiURL: "https://en.wikipedia.org/wiki/Lucy_Lawless"},
		{Name: "Sam Walton", FamousFor: "Founder of Walmart and Sam's Club", WikiURL: "https://en.wikipedia.org/wiki/Sam_Walton"},
	},
	"Mar-30": {
		{Name: "Vincent van Gogh", FamousFor: "Post-impressionist painter ('The Starry Night')", WikiURL: "https://en.wikipedia.org/wiki/Vincent_van_Gogh"},
		{Name: "Celine Dion", FamousFor: "Grammy-winning singer ('My Heart Will Go On')", WikiURL: "https://en.wikipedia.org/wiki/Celine_Dion"},
		{Name: "Eric Clapton", FamousFor: "Rock and blues guitarist and singer", WikiURL: "https://en.wikipedia.org/wiki/Eric_Clapton"},
	},
	"Mar-31": {
		{Name: "René Descartes", FamousFor: "Philosopher ('I think, therefore I am')", WikiURL: "https://en.wikipedia.org/wiki/Ren%C3%A9_Descartes"},
		{Name: "Ewan McGregor", FamousFor: "Actor ('Obi-Wan Kenobi' in 'Star Wars')", WikiURL: "https://en.wikipedia.org/wiki/Ewan_McGregor"},
		{Name: "Christopher Walken", FamousFor: "Oscar-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Christopher_Walken"},
	},
	"Apr-01": {
		{Name: "Lon Chaney", FamousFor: "Silent film actor, 'The Man of a Thousand Faces'", WikiURL: "https://en.wikipedia.org/wiki/Lon_Chaney"},
		{Name: "Rachel Maddow", FamousFor: "Political commentator and television host", WikiURL: "https://en.wikipedia.org/wiki/Rachel_Maddow"},
		{Name: "Susan Boyle", FamousFor: "Singer who rose to fame on 'Britain's Got Talent'", WikiURL: "https://en.wikipedia.org/wiki/Susan_Boyle"},
	},
	"Apr-02": {
		{Name: "Hans Christian Andersen", FamousFor: "Author of fairy tales ('The Little Mermaid')", WikiURL: "https://en.wikipedia.org/wiki/Hans_Christian_Andersen"},
		{Name: "Marvin Gaye", FamousFor: "Motown singer and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Marvin_Gaye"},
		{Name: "Michael Fassbender", FamousFor: "Actor ('X-Men', 'Steve Jobs')", WikiURL: "https://en.wikipedia.org/wiki/Michael_Fassbender"},
	},
	"Apr-03": {
		{Name: "Marlon Brando", FamousFor: "Oscar-winning actor ('The Godfather', 'On the Waterfront')", WikiURL: "https://en.wikipedia.org/wiki/Marlon_Brando"},
		{Name: "Eddie Murphy", FamousFor: "Comedian and actor", WikiURL: "https://en.wikipedia.org/wiki/Eddie_Murphy"},
		{Name: "Alec Baldwin", FamousFor: "Emmy-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Alec_Baldwin"},
	},
	"Apr-04": {
		{Name: "Maya Angelou", FamousFor: "Poet, memoirist, and civil rights activist", WikiURL: "https://en.wikipedia.org/wiki/Maya_Angelou"},
		{Name: "Robert Downey Jr.", FamousFor: "Actor ('Iron Man' in Marvel films)", WikiURL: "https://en.wikipedia.org/wiki/Robert_Downey_Jr."},
		{Name: "Heath Ledger", FamousFor: "Oscar-winning actor ('The Joker' in 'The Dark Knight')", WikiURL: "https://en.wikipedia.org/wiki/Heath_Ledger"},
	},
	"Apr-05": {
		{Name: "Booker T. Washington", FamousFor: "Educator, author, and advisor to presidents", WikiURL: "https://en.wikipedia.org/wiki/Booker_T._Washington"},
		{Name: "Pharrell Williams", FamousFor: "Grammy-winning musician, producer, and fashion designer", WikiURL: "https://en.wikipedia.org/wiki/Pharrell_Williams"},
		{Name: "Colin Powell", FamousFor: "U.S. general and Secretary of State", WikiURL: "https://en.wikipedia.org/wiki/Colin_Powell"},
	},
	"Apr-06": {
		{Name: "Raphael", FamousFor: "High Renaissance painter and architect", WikiURL: "https://en.wikipedia.org/wiki/Raphael"},
		{Name: "Paul Rudd", FamousFor: "Actor ('Ant-Man' in Marvel films)", WikiURL: "https://en.wikipedia.org/wiki/Paul_Rudd"},
		{Name: "Zach Braff", FamousFor: "Actor ('J.D.' in 'Scrubs')", WikiURL: "https://en.wikipedia.org/wiki/Zach_Braff"},
	},
	"Apr-07": {
		{Name: "Billie Holiday", FamousFor: "Pioneering jazz and swing singer", WikiURL: "https://en.wikipedia.org/wiki/Billie_Holiday"},
		{Name: "Francis Ford Coppola", FamousFor: "Oscar-winning director ('The Godfather', 'Apocalypse Now')", WikiURL: "https://en.wikipedia.org/wiki/Francis_Ford_Coppola"},
		{Name: "Jackie Chan", FamousFor: "Martial artist and actor", WikiURL: "https://en.wikipedia.org/wiki/Jackie_Chan"},
	},
	"Apr-08": {
		{Name: "Buddha (Siddhartha Gautama)", FamousFor: "Founder of Buddhism (traditional birthday)", WikiURL: "https://en.wikipedia.org/wiki/The_Buddha"},
		{Name: "Robin Wright", FamousFor: "Actress ('Forrest Gump', 'House of Cards')", WikiURL: "https://en.wikipedia.org/wiki/Robin_Wright"},
		{Name: "Patricia Arquette", FamousFor: "Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Patricia_Arquette"},
	},
	"Apr-09": {
		{Name: "Hugh Hefner", FamousFor: "Founder of Playboy Magazine", WikiURL: "https://en.wikipedia.org/wiki/Hugh_Hefner"},
		{Name: "Kristen Stewart", FamousFor: "Actress ('Twilight' series)", WikiURL: "https://en.wikipedia.org/wiki/Kristen_Stewart"},
		{Name: "Dennis Quaid", FamousFor: "Actor ('The Parent Trap', 'The Day After Tomorrow')", WikiURL: "https://en.wikipedia.org/wiki/Dennis_Quaid"},
	},
	"Apr-10": {
		{Name: "Joseph Pulitzer", FamousFor: "Newspaper publisher, namesake of the Pulitzer Prize", WikiURL: "https://en.wikipedia.org/wiki/Joseph_Pulitzer"},
		{Name: "Steven Seagal", FamousFor: "Action film actor and martial artist", WikiURL: "https://en.wikipedia.org/wiki/Steven_Seagal"},
		{Name: "Mandy Moore", FamousFor: "Singer and actress ('This Is Us')", WikiURL: "https://en.wikipedia.org/wiki/Mandy_Moore"},
	},
	"Apr-11": {
		{Name: "Bill Irwin", FamousFor: "Tony-winning actor and clown", WikiURL: "https://en.wikipedia.org/wiki/Bill_Irwin"},
		{Name: "Jeremy Clarkson", FamousFor: "Television host ('Top Gear', 'The Grand Tour')", WikiURL: "https://en.wikipedia.org/wiki/Jeremy_Clarkson"},
		{Name: "Joss Stone", FamousFor: "Grammy-winning singer and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Joss_Stone"},
	},
	"Apr-12": {
		{Name: "David Letterman", FamousFor: "Long-time late-night television host", WikiURL: "https://en.wikipedia.org/wiki/David_Letterman"},
		{Name: "Saoirse Ronan", FamousFor: "Actress ('Lady Bird', 'Little Women')", WikiURL: "https://en.wikipedia.org/wiki/Saoirse_Ronan"},
		{Name: "Claire Danes", FamousFor: "Emmy-winning actress ('Homeland')", WikiURL: "https://en.wikipedia.org/wiki/Claire_Danes"},
	},
	"Apr-13": {
		{Name: "Thomas Jefferson", FamousFor: "3rd U.S. President, author of the Declaration of Independence", WikiURL: "https://en.wikipedia.org/wiki/Thomas_Jefferson"},
		{Name: "Al Green", FamousFor: "Grammy-winning soul singer", WikiURL: "https://en.wikipedia.org/wiki/Al_Green"},
		{Name: "Ron Perlman", FamousFor: "Actor ('Hellboy', 'Sons of Anarchy')", WikiURL: "https://en.wikipedia.org/wiki/Ron_Perlman"},
	},
	"Apr-14": {
		{Name: "Adrien Brody", FamousFor: "Oscar-winning actor ('The Pianist')", WikiURL: "https://en.wikipedia.org/wiki/Adrien_Brody"},
		{Name: "Sarah Michelle Gellar", FamousFor: "Actress ('Buffy the Vampire Slayer')", WikiURL: "https://en.wikipedia.org/wiki/Sarah_Michelle_Gellar"},
		{Name: "Pete Rose", FamousFor: "Professional baseball player and manager", WikiURL: "https://en.wikipedia.org/wiki/Pete_Rose"},
	},
	"Apr-15": {
		{Name: "Leonardo da Vinci", FamousFor: "High Renaissance artist and inventor ('Mona Lisa')", WikiURL: "https://en.wikipedia.org/wiki/Leonardo_da_Vinci"},
		{Name: "Emma Watson", FamousFor: "Actress ('Hermione Granger' in 'Harry Potter')", WikiURL: "https://en.wikipedia.org/wiki/Emma_Watson"},
		{Name: "Seth Rogen", FamousFor: "Actor, comedian, and filmmaker", WikiURL: "https://en.wikipedia.org/wiki/Seth_Rogen"},
	},
	"Apr-16": {
		{Name: "Charlie Chaplin", FamousFor: "Iconic silent film actor and comedian", WikiURL: "https://en.wikipedia.org/wiki/Charlie_Chaplin"},
		{Name: "Pope Benedict XVI", FamousFor: "Head of the Catholic Church from 2005 to 2013", WikiURL: "https://en.wikipedia.org/wiki/Pope_Benedict_XVI"},
		{Name: "Martin Lawrence", FamousFor: "Comedian and actor ('Bad Boys')", WikiURL: "https://en.wikipedia.org/wiki/Martin_Lawrence"},
	},
	"Apr-17": {
		{Name: "J.P. Morgan", FamousFor: "Financier and banker", WikiURL: "https://en.wikipedia.org/wiki/J._P._Morgan"},
		{Name: "Victoria Beckham", FamousFor: "Singer (Spice Girls) and fashion designer", WikiURL: "https://en.wikipedia.org/wiki/Victoria_Beckham"},
		{Name: "Jennifer Garner", FamousFor: "Actress ('Alias', '13 Going on 30')", WikiURL: "https://en.wikipedia.org/wiki/Jennifer_Garner"},
	},
	"Apr-18": {
		{Name: "Conan O'Brien", FamousFor: "Late-night television host", WikiURL: "https://en.wikipedia.org/wiki/Conan_O%27Brien"},
		{Name: "Kourtney Kardashian", FamousFor: "Media personality and socialite", WikiURL: "https://en.wikipedia.org/wiki/Kourtney_Kardashian"},
		{Name: "James Woods", FamousFor: "Emmy-winning actor", WikiURL: "https://en.wikipedia.org/wiki/James_Woods"},
	},
	"Apr-19": {
		{Name: "Kate Hudson", FamousFor: "Actress ('Almost Famous')", WikiURL: "https://en.wikipedia.org/wiki/Kate_Hudson"},
		{Name: "James Franco", FamousFor: "Actor and filmmaker", WikiURL: "https://en.wikipedia.org/wiki/James_Franco"},
		{Name: "Maria Sharapova", FamousFor: "Professional tennis player", WikiURL: "https://en.wikipedia.org/wiki/Maria_Sharapova"},
	},
	"Apr-20": {
		{Name: "Adolf Hitler", FamousFor: "Dictator of Nazi Germany", WikiURL: "https://en.wikipedia.org/wiki/Adolf_Hitler"},
		{Name: "Carmen Electra", FamousFor: "Model and actress ('Baywatch')", WikiURL: "https://en.wikipedia.org/wiki/Carmen_Electra"},
		{Name: "Jessica Lange", FamousFor: "Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Jessica_Lange"},
	},
	"Apr-21": {
		{Name: "Queen Elizabeth II", FamousFor: "Queen of the United Kingdom from 1952 to 2022", WikiURL: "https://en.wikipedia.org/wiki/Elizabeth_II"},
		{Name: "Iggy Pop", FamousFor: "Musician, 'Godfather of Punk'", WikiURL: "https://en.wikipedia.org/wiki/Iggy_Pop"},
		{Name: "James McAvoy", FamousFor: "Actor ('Professor X' in 'X-Men')", WikiURL: "https://en.wikipedia.org/wiki/James_McAvoy"},
	},
	"Apr-22": {
		{Name: "Vladimir Lenin", FamousFor: "Russian revolutionary and head of Soviet Russia", WikiURL: "https://en.wikipedia.org/wiki/Vladimir_Lenin"},
		{Name: "Jack Nicholson", FamousFor: "Three-time Oscar-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Jack_Nicholson"},
		{Name: "Amber Heard", FamousFor: "Actress", WikiURL: "https://en.wikipedia.org/wiki/Amber_Heard"},
	},
	"Apr-23": {
		{Name: "William Shakespeare", FamousFor: "Playwright and poet, widely regarded as the greatest writer in English", WikiURL: "https://en.wikipedia.org/wiki/William_Shakespeare"},
		{Name: "Shirley Temple", FamousFor: "Child actress and diplomat", WikiURL: "https://en.wikipedia.org/wiki/Shirley_Temple"},
		{Name: "John Cena", FamousFor: "Professional wrestler and actor", WikiURL: "https://en.wikipedia.org/wiki/John_Cena"},
	},
	"Apr-24": {
		{Name: "Barbra Streisand", FamousFor: "EGOT-winning singer and actress", WikiURL: "https://en.wikipedia.org/wiki/Barbra_Streisand"},
		{Name: "Kelly Clarkson", FamousFor: "Grammy-winning singer, original 'American Idol'", WikiURL: "https://en.wikipedia.org/wiki/Kelly_Clarkson"},
		{Name: "Cedric the Entertainer", FamousFor: "Comedian and actor", WikiURL: "https://en.wikipedia.org/wiki/Cedric_the_Entertainer"},
	},
	"Apr-25": {
		{Name: "Al Pacino", FamousFor: "Oscar-winning actor ('The Godfather', 'Scent of a Woman')", WikiURL: "https://en.wikipedia.org/wiki/Al_Pacino"},
		{Name: "Renée Zellweger", FamousFor: "Oscar-winning actress ('Bridget Jones's Diary', 'Judy')", WikiURL: "https://en.wikipedia.org/wiki/Ren%C3%A9e_Zellweger"},
		{Name: "Ella Fitzgerald", FamousFor: "Jazz singer, 'First Lady of Song'", WikiURL: "https://en.wikipedia.org/wiki/Ella_Fitzgerald"},
	},
	"Apr-26": {
		{Name: "Jet Li", FamousFor: "Martial artist and actor", WikiURL: "https://en.wikipedia.org/wiki/Jet_Li"},
		{Name: "Channing Tatum", FamousFor: "Actor ('Magic Mike')", WikiURL: "https://en.wikipedia.org/wiki/Channing_Tatum"},
		{Name: "Kevin James", FamousFor: "Comedian and actor ('The King of Queens')", WikiURL: "https://en.wikipedia.org/wiki/Kevin_James"},
	},
	"Apr-27": {
		{Name: "Ulysses S. Grant", FamousFor: "18th U.S. President and Union Army General", WikiURL: "https://en.wikipedia.org/wiki/Ulysses_S._Grant"},
		{Name: "Lizzo", FamousFor: "Grammy-winning singer and rapper", WikiURL: "https://en.wikipedia.org/wiki/Lizzo"},
		{Name: "Sheena Easton", FamousFor: "Grammy-winning singer", WikiURL: "https://en.wikipedia.org/wiki/Sheena_Easton"},
	},
	"Apr-28": {
		{Name: "Harper Lee", FamousFor: "Author ('To Kill a Mockingbird')", WikiURL: "https://en.wikipedia.org/wiki/Harper_Lee"},
		{Name: "Jay Leno", FamousFor: "Former host of 'The Tonight Show'", WikiURL: "https://en.wikipedia.org/wiki/Jay_Leno"},
		{Name: "Penélope Cruz", FamousFor: "Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Pen%C3%A9lope_Cruz"},
	},
	"Apr-29": {
		{Name: "Daniel Day-Lewis", FamousFor: "Three-time Oscar-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Daniel_Day-Lewis"},
		{Name: "Michelle Pfeiffer", FamousFor: "Actress ('Batman Returns', 'Scarface')", WikiURL: "https://en.wikipedia.org/wiki/Michelle_Pfeiffer"},
		{Name: "Uma Thurman", FamousFor: "Actress ('Pulp Fiction', 'Kill Bill')", WikiURL: "https://en.wikipedia.org/wiki/Uma_Thurman"},
	},
	"Apr-30": {
		{Name: "Kirsten Dunst", FamousFor: "Actress ('Spider-Man', 'Bring It On')", WikiURL: "https://en.wikipedia.org/wiki/Kirsten_Dunst"},
		{Name: "Gal Gadot", FamousFor: "Actress ('Wonder Woman')", WikiURL: "https://en.wikipedia.org/wiki/Gal_Gadot"},
		{Name: "Willie Nelson", FamousFor: "Country music singer-songwriter", WikiURL: "https://en.wikipedia.org/wiki/Willie_Nelson"},
	},
	"May-01": {
		{Name: "Wes Anderson", FamousFor: "Film director with a distinctive visual style", WikiURL: "https://en.wikipedia.org/wiki/Wes_Anderson"},
		{Name: "Tim McGraw", FamousFor: "Country music singer", WikiURL: "https://en.wikipedia.org/wiki/Tim_McGraw"},
		{Name: "Jamie Dornan", FamousFor: "Actor ('Fifty Shades' series)", WikiURL: "https://en.wikipedia.org/wiki/Jamie_Dornan"},
	},
	"May-02": {
		{Name: "Dwayne 'The Rock' Johnson", FamousFor: "Professional wrestler and actor", WikiURL: "https://en.wikipedia.org/wiki/Dwayne_Johnson"},
		{Name: "David Beckham", FamousFor: "Professional footballer and style icon", WikiURL: "https://en.wikipedia.org/wiki/David_Beckham"},
		{Name: "Lily Allen", FamousFor: "Singer and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Lily_Allen"},
	},
	"May-03": {
		{Name: "James Brown", FamousFor: "Singer, 'The Godfather of Soul'", WikiURL: "https://en.wikipedia.org/wiki/James_Brown"},
		{Name: "Christina Hendricks", FamousFor: "Actress ('Mad Men')", WikiURL: "https://en.wikipedia.org/wiki/Christina_Hendricks"},
		{Name: "Bing Crosby", FamousFor: "Singer and actor ('White Christmas')", WikiURL: "https://en.wikipedia.org/wiki/Bing_Crosby"},
	},
	"May-04": {
		{Name: "Audrey Hepburn", FamousFor: "Actress, fashion icon, and humanitarian ('Breakfast at Tiffany's')", WikiURL: "https://en.wikipedia.org/wiki/Audrey_Hepburn"},
		{Name: "Will Arnett", FamousFor: "Actor and comedian ('Arrested Development', 'BoJack Horseman')", WikiURL: "https://en.wikipedia.org/wiki/Will_Arnett"},
		{Name: "Erin Andrews", FamousFor: "Sportscaster and television personality", WikiURL: "https://en.wikipedia.org/wiki/Erin_Andrews"},
	},
	"May-05": {
		{Name: "Karl Marx", FamousFor: "Philosopher, economist, and revolutionary socialist", WikiURL: "https://en.wikipedia.org/wiki/Karl_Marx"},
		{Name: "Adele", FamousFor: "Grammy-winning singer and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Adele"},
		{Name: "Henry Cavill", FamousFor: "Actor ('Superman', 'The Witcher')", WikiURL: "https://en.wikipedia.org/wiki/Henry_Cavill"},
	},
	"May-06": {
		{Name: "Sigmund Freud", FamousFor: "Founder of psychoanalysis", WikiURL: "https://en.wikipedia.org/wiki/Sigmund_Freud"},
		{Name: "George Clooney", FamousFor: "Oscar-winning actor and director", WikiURL: "https://en.wikipedia.org/wiki/George_Clooney"},
		{Name: "Orson Welles", FamousFor: "Director, writer, and actor ('Citizen Kane')", WikiURL: "https://en.wikipedia.org/wiki/Orson_Welles"},
	},
	"May-07": {
		{Name: "Pyotr Ilyich Tchaikovsky", FamousFor: "Composer ('Swan Lake', 'The Nutcracker')", WikiURL: "https://en.wikipedia.org/wiki/Pyotr_Ilyich_Tchaikovsky"},
		{Name: "Eva Perón", FamousFor: "First Lady of Argentina", WikiURL: "https://en.wikipedia.org/wiki/Eva_Per%C3%B3n"},
		{Name: "Gary Cooper", FamousFor: "Oscar-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Gary_Cooper"},
	},
	"May-08": {
		{Name: "Harry S. Truman", FamousFor: "33rd U.S. President", WikiURL: "https://en.wikipedia.org/wiki/Harry_S._Truman"},
		{Name: "David Attenborough", FamousFor: "Broadcaster and natural historian", WikiURL: "https://en.wikipedia.org/wiki/David_Attenborough"},
		{Name: "Don Rickles", FamousFor: "Comedian and actor, master of insult comedy", WikiURL: "https://en.wikipedia.org/wiki/Don_Rickles"},
	},
	"May-09": {
		{Name: "Billy Joel", FamousFor: "Singer-songwriter and pianist ('Piano Man')", WikiURL: "https://en.wikipedia.org/wiki/Billy_Joel"},
		{Name: "Rosario Dawson", FamousFor: "Actress ('Ahsoka')", WikiURL: "https://en.wikipedia.org/wiki/Rosario_Dawson"},
		{Name: "Candice Bergen", FamousFor: "Emmy-winning actress ('Murphy Brown')", WikiURL: "https://en.wikipedia.org/wiki/Candice_Bergen"},
	},
	"May-10": {
		{Name: "Bono", FamousFor: "Lead singer of U2 and activist", WikiURL: "https://en.wikipedia.org/wiki/Bono"},
		{Name: "Fred Astaire", FamousFor: "Dancer and actor", WikiURL: "https://en.wikipedia.org/wiki/Fred_Astaire"},
		{Name: "John Wilkes Booth", FamousFor: "Assassin of Abraham Lincoln", WikiURL: "https://en.wikipedia.org/wiki/John_Wilkes_Booth"},
	},
	"May-11": {
		{Name: "Salvador Dalí", FamousFor: "Surrealist artist", WikiURL: "https://en.wikipedia.org/wiki/Salvador_Dal%C3%AD"},
		{Name: "Cam Newton", FamousFor: "NFL MVP quarterback", WikiURL: "https://en.wikipedia.org/wiki/Cam_Newton"},
		{Name: "Cory Monteith", FamousFor: "Actor ('Finn Hudson' in 'Glee')", WikiURL: "https://en.wikipedia.org/wiki/Cory_Monteith"},
	},
	"May-12": {
		{Name: "Florence Nightingale", FamousFor: "Founder of modern nursing", WikiURL: "https://en.wikipedia.org/wiki/Florence_Nightingale"},
		{Name: "Tony Hawk", FamousFor: "Professional skateboarder and entrepreneur", WikiURL: "https://en.wikipedia.org/wiki/Tony_Hawk"},
		{Name: "Rami Malek", FamousFor: "Oscar-winning actor ('Bohemian Rhapsody')", WikiURL: "https://en.wikipedia.org/wiki/Rami_Malek"},
	},
	"May-13": {
		{Name: "Stevie Wonder", FamousFor: "Grammy-winning musician and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Stevie_Wonder"},
		{Name: "Robert Pattinson", FamousFor: "Actor ('Twilight', 'The Batman')", WikiURL: "https://en.wikipedia.org/wiki/Robert_Pattinson"},
		{Name: "Stephen Colbert", FamousFor: "Comedian and late-night talk show host", WikiURL: "https://en.wikipedia.org/wiki/Stephen_Colbert"},
	},
	"May-14": {
		{Name: "George Lucas", FamousFor: "Creator of 'Star Wars' and 'Indiana Jones'", WikiURL: "https://en.wikipedia.org/wiki/George_Lucas"},
		{Name: "Mark Zuckerberg", FamousFor: "Co-founder of Facebook", WikiURL: "https://en.wikipedia.org/wiki/Mark_Zuckerberg"},
		{Name: "Cate Blanchett", FamousFor: "Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Cate_Blanchett"},
	},
	"May-15": {
		{Name: "Andy Murray", FamousFor: "Professional tennis player, Wimbledon champion", WikiURL: "https://en.wikipedia.org/wiki/Andy_Murray"},
		{Name: "L. Frank Baum", FamousFor: "Author of 'The Wonderful Wizard of Oz'", WikiURL: "https://en.wikipedia.org/wiki/L._Frank_Baum"},
		{Name: "Zara Tindall", FamousFor: "British Royal Family member and equestrian", WikiURL: "https://en.wikipedia.org/wiki/Zara_Tindall"},
	},
	"May-16": {
		{Name: "Megan Fox", FamousFor: "Actress ('Transformers')", WikiURL: "https://en.wikipedia.org/wiki/Megan_Fox"},
		{Name: "Janet Jackson", FamousFor: "Grammy-winning singer and dancer", WikiURL: "https://en.wikipedia.org/wiki/Janet_Jackson"},
		{Name: "Pierce Brosnan", FamousFor: "Actor ('James Bond')", WikiURL: "https://en.wikipedia.org/wiki/Pierce_Brosnan"},
	},
	"May-17": {
		{Name: "Sugar Ray Leonard", FamousFor: "Olympic gold medalist professional boxer", WikiURL: "https://en.wikipedia.org/wiki/Sugar_Ray_Leonard"},
		{Name: "Dennis Hopper", FamousFor: "Actor and filmmaker ('Easy Rider')", WikiURL: "https://en.wikipedia.org/wiki/Dennis_Hopper"},
		{Name: "Trent Reznor", FamousFor: "Frontman of Nine Inch Nails, Oscar-winning composer", WikiURL: "https://en.wikipedia.org/wiki/Trent_Reznor"},
	},
	"May-18": {
		{Name: "Tina Fey", FamousFor: "Emmy-winning writer, actress, and comedian ('30 Rock')", WikiURL: "https://en.wikipedia.org/wiki/Tina_Fey"},
		{Name: "Pope John Paul II", FamousFor: "Head of the Catholic Church from 1978 to 2005", WikiURL: "https://en.wikipedia.org/wiki/Pope_John_Paul_II"},
		{Name: "Jack Johnson", FamousFor: "Singer-songwriter and musician", WikiURL: "https://en.wikipedia.org/wiki/Jack_Johnson_(musician)"},
	},
	"May-19": {
		{Name: "Malcolm X", FamousFor: "Human rights activist", WikiURL: "https://en.wikipedia.org/wiki/Malcolm_X"},
		{Name: "Sam Smith", FamousFor: "Grammy-winning singer and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Sam_Smith"},
		{Name: "Pete Townshend", FamousFor: "Guitarist and songwriter for The Who", WikiURL: "https://en.wikipedia.org/wiki/Pete_Townshend"},
	},
	"May-20": {
		{Name: "Cher", FamousFor: "EGOT-winning singer and actress, 'Goddess of Pop'", WikiURL: "https://en.wikipedia.org/wiki/Cher"},
		{Name: "James Stewart", FamousFor: "Oscar-winning actor ('It's a Wonderful Life')", WikiURL: "https://en.wikipedia.org/wiki/James_Stewart"},
		{Name: "Busta Rhymes", FamousFor: "Rapper and actor", WikiURL: "https://en.wikipedia.org/wiki/Busta_Rhymes"},
	},
	"May-21": {
		{Name: "The Notorious B.I.G.", FamousFor: "Influential rapper", WikiURL: "https://en.wikipedia.org/wiki/The_Notorious_B.I.G."},
		{Name: "Mr. T", FamousFor: "Actor ('The A-Team')", WikiURL: "https://en.wikipedia.org/wiki/Mr._T"},
		{Name: "Lisa Edelstein", FamousFor: "Actress ('Dr. Lisa Cuddy' in 'House')", WikiURL: "https://en.wikipedia.org/wiki/Lisa_Edelstein"},
	},
	"May-22": {
		{Name: "Arthur Conan Doyle", FamousFor: "Author, creator of Sherlock Holmes", WikiURL: "https://en.wikipedia.org/wiki/Arthur_Conan_Doyle"},
		{Name: "Naomi Campbell", FamousFor: "Supermodel", WikiURL: "https://en.wikipedia.org/wiki/Naomi_Campbell"},
		{Name: "Morrissey", FamousFor: "Lead singer of The Smiths", WikiURL: "https://en.wikipedia.org/wiki/Morrissey"},
	},
	"May-23": {
		{Name: "Drew Carey", FamousFor: "Comedian and television host", WikiURL: "https://en.wikipedia.org/wiki/Drew_Carey"},
		{Name: "Joan Collins", FamousFor: "Actress ('Dynasty')", WikiURL: "https://en.wikipedia.org/wiki/Joan_Collins"},
		{Name: "Jewel", FamousFor: "Singer-songwriter", WikiURL: "https://en.wikipedia.org/wiki/Jewel_(singer)"},
	},
	"May-24": {
		{Name: "Bob Dylan", FamousFor: "Nobel Prize-winning singer-songwriter", WikiURL: "https://en.wikipedia.org/wiki/Bob_Dylan"},
		{Name: "Queen Victoria", FamousFor: "Queen of the United Kingdom from 1837 to 1901", WikiURL: "https://en.wikipedia.org/wiki/Queen_Victoria"},
		{Name: "John C. Reilly", FamousFor: "Actor and comedian", WikiURL: "https://en.wikipedia.org/wiki/John_C._Reilly"},
	},
	"May-25": {
		{Name: "Cillian Murphy", FamousFor: "Oscar-winning actor ('Oppenheimer', 'Peaky Blinders')", WikiURL: "https://en.wikipedia.org/wiki/Cillian_Murphy"},
		{Name: "Mike Myers", FamousFor: "Comedian and actor ('Austin Powers', 'Shrek')", WikiURL: "https://en.wikipedia.org/wiki/Mike_Myers"},
		{Name: "Ian McKellen", FamousFor: "Actor ('Gandalf' in 'Lord of the Rings', 'Magneto' in 'X-Men')", WikiURL: "https://en.wikipedia.org/wiki/Ian_McKellen"},
	},
	"May-26": {
		{Name: "Miles Davis", FamousFor: "Influential jazz trumpeter and composer", WikiURL: "https://en.wikipedia.org/wiki/Miles_Davis"},
		{Name: "John Wayne", FamousFor: "Actor, icon of Western films", WikiURL: "https://en.wikipedia.org/wiki/John_Wayne"},
		{Name: "Helena Bonham Carter", FamousFor: "Actress", WikiURL: "https://en.wikipedia.org/wiki/Helena_Bonham_Carter"},
	},
	"May-27": {
		{Name: "Vincent Price", FamousFor: "Actor known for his roles in horror films", WikiURL: "https://en.wikipedia.org/wiki/Vincent_Price"},
		{Name: "Christopher Lee", FamousFor: "Actor ('Dracula', 'Saruman' in 'Lord of the Rings')", WikiURL: "https://en.wikipedia.org/wiki/Christopher_Lee"},
		{Name: "Jamie Oliver", FamousFor: "Chef and television personality", WikiURL: "https://en.wikipedia.org/wiki/Jamie_Oliver"},
	},
	"May-28": {
		{Name: "Kylie Minogue", FamousFor: "Singer, 'Princess of Pop'", WikiURL: "https://en.wikipedia.org/wiki/Kylie_Minogue"},
		{Name: "Carey Mulligan", FamousFor: "Actress", WikiURL: "https://en.wikipedia.org/wiki/Carey_Mulligan"},
		{Name: "Ian Fleming", FamousFor: "Author, creator of James Bond", WikiURL: "https://en.wikipedia.org/wiki/Ian_Fleming"},
	},
	"May-29": {
		{Name: "John F. Kennedy", FamousFor: "35th U.S. President", WikiURL: "https://en.wikipedia.org/wiki/John_F._Kennedy"},
		{Name: "Annette Bening", FamousFor: "Actress", WikiURL: "https://en.wikipedia.org/wiki/Annette_Bening"},
		{Name: "Mel B", FamousFor: "Singer ('Scary Spice' of the Spice Girls)", WikiURL: "https://en.wikipedia.org/wiki/Mel_B"},
	},
	"May-30": {
		{Name: "Idina Menzel", FamousFor: "Tony-winning actress and singer ('Elsa' in 'Frozen')", WikiURL: "https://en.wikipedia.org/wiki/Idina_Menzel"},
		{Name: "CeeLo Green", FamousFor: "Singer and songwriter", WikiURL: "https://en.wikipedia.org/wiki/CeeLo_Green"},
		{Name: "Manny Ramirez", FamousFor: "Professional baseball player", WikiURL: "https://en.wikipedia.org/wiki/Manny_Ramirez"},
	},
	"May-31": {
		{Name: "Clint Eastwood", FamousFor: "Oscar-winning actor and director", WikiURL: "https://en.wikipedia.org/wiki/Clint_Eastwood"},
		{Name: "Colin Farrell", FamousFor: "Actor", WikiURL: "https://en.wikipedia.org/wiki/Colin_Farrell"},
		{Name: "Brooke Shields", FamousFor: "Actress and model", WikiURL: "https://en.wikipedia.org/wiki/Brooke_Shields"},
	},
	"Jun-01": {
		{Name: "Marilyn Monroe", FamousFor: "Actress, singer, and model; an icon of pop culture", WikiURL: "https://en.wikipedia.org/wiki/Marilyn_Monroe"},
		{Name: "Morgan Freeman", FamousFor: "Oscar-winning actor and narrator", WikiURL: "https://en.wikipedia.org/wiki/Morgan_Freeman"},
		{Name: "Heidi Klum", FamousFor: "Supermodel and television host", WikiURL: "https://en.wikipedia.org/wiki/Heidi_Klum"},
	},
	"Jun-02": {
		{Name: "Zachary Quinto", FamousFor: "Actor ('Spock' in 'Star Trek' reboot series)", WikiURL: "https://en.wikipedia.org/wiki/Zachary_Quinto"},
		{Name: "Justin Long", FamousFor: "Actor ('Dodgeball', 'Jeepers Creepers')", WikiURL: "https://en.wikipedia.org/wiki/Justin_Long"},
		{Name: "Wayne Brady", FamousFor: "Comedian, actor, and singer ('Whose Line Is It Anyway?')", WikiURL: "https://en.wikipedia.org/wiki/Wayne_Brady"},
	},
	"Jun-03": {
		{Name: "Rafael Nadal", FamousFor: "Professional tennis player, multiple Grand Slam winner", WikiURL: "https://en.wikipedia.org/wiki/Rafael_Nadal"},
		{Name: "Anderson Cooper", FamousFor: "Journalist and television anchor", WikiURL: "https://en.wikipedia.org/wiki/Anderson_Cooper"},
		{Name: "Jodie Whittaker", FamousFor: "Actress (Thirteenth Doctor in 'Doctor Who')", WikiURL: "https://en.wikipedia.org/wiki/Jodie_Whittaker"},
	},
	"Jun-04": {
		{Name: "Angelina Jolie", FamousFor: "Oscar-winning actress and humanitarian", WikiURL: "https://en.wikipedia.org/wiki/Angelina_Jolie"},
		{Name: "Russell Brand", FamousFor: "Comedian and actor", WikiURL: "https://en.wikipedia.org/wiki/Russell_Brand"},
		{Name: "Bar Refaeli", FamousFor: "Model and television host", WikiURL: "https://en.wikipedia.org/wiki/Bar_Refaeli"},
	},
	"Jun-05": {
		{Name: "Mark Wahlberg", FamousFor: "Actor and producer", WikiURL: "https://en.wikipedia.org/wiki/Mark_Wahlberg"},
		{Name: "Kenny G", FamousFor: "Grammy-winning smooth jazz saxophonist", WikiURL: "https://en.wikipedia.org/wiki/Kenny_G"},
		{Name: "Pete Wentz", FamousFor: "Bassist and lyricist for Fall Out Boy", WikiURL: "https://en.wikipedia.org/wiki/Pete_Wentz"},
	},
	"Jun-06": {
		{Name: "Paul Giamatti", FamousFor: "Emmy-winning actor ('Sideways', 'John Adams')", WikiURL: "https://en.wikipedia.org/wiki/Paul_Giamatti"},
		{Name: "Jason Isaacs", FamousFor: "Actor ('Lucius Malfoy' in 'Harry Potter')", WikiURL: "https://en.wikipedia.org/wiki/Jason_Isaacs"},
		{Name: "Robert Englund", FamousFor: "Actor ('Freddy Krueger' in 'A Nightmare on Elm Street')", WikiURL: "https://en.wikipedia.org/wiki/Robert_Englund"},
	},
	"Jun-07": {
		{Name: "Prince", FamousFor: "Iconic, Grammy-winning musician and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Prince_(musician)"},
		{Name: "Liam Neeson", FamousFor: "Actor ('Schindler's List', 'Taken')", WikiURL: "https://en.wikipedia.org/wiki/Liam_Neeson"},
		{Name: "Michael Cera", FamousFor: "Actor ('Superbad', 'Juno')", WikiURL: "https://en.wikipedia.org/wiki/Michael_Cera"},
	},
	"Jun-08": {
		{Name: "Kanye West", FamousFor: "Grammy-winning rapper, producer, and fashion designer", WikiURL: "https://en.wikipedia.org/wiki/Kanye_West"},
		{Name: "Barbara Bush", FamousFor: "Former First Lady of the United States", WikiURL: "https://en.wikipedia.org/wiki/Barbara_Bush"},
		{Name: "Frank Lloyd Wright", FamousFor: "Influential architect", WikiURL: "https://en.wikipedia.org/wiki/Frank_Lloyd_Wright"},
	},
	"Jun-09": {
		{Name: "Johnny Depp", FamousFor: "Actor ('Captain Jack Sparrow' in 'Pirates of the Caribbean')", WikiURL: "https://en.wikipedia.org/wiki/Johnny_Depp"},
		{Name: "Natalie Portman", FamousFor: "Oscar-winning actress ('Black Swan')", WikiURL: "https://en.wikipedia.org/wiki/Natalie_Portman"},
		{Name: "Michael J. Fox", FamousFor: "Actor ('Marty McFly' in 'Back to the Future')", WikiURL: "https://en.wikipedia.org/wiki/Michael_J._Fox"},
	},
	"Jun-10": {
		{Name: "Judy Garland", FamousFor: "Actress and singer ('Dorothy' in 'The Wizard of Oz')", WikiURL: "https://en.wikipedia.org/wiki/Judy_Garland"},
		{Name: "Prince Philip, Duke of Edinburgh", FamousFor: "Husband of Queen Elizabeth II", WikiURL: "https://en.wikipedia.org/wiki/Prince_Philip,_Duke_of_Edinburgh"},
		{Name: "Kate Upton", FamousFor: "Model and actress", WikiURL: "https://en.wikipedia.org/wiki/Kate_Upton"},
	},
	"Jun-11": {
		{Name: "Gene Wilder", FamousFor: "Actor ('Willy Wonka & the Chocolate Factory')", WikiURL: "https://en.wikipedia.org/wiki/Gene_Wilder"},
		{Name: "Shia LaBeouf", FamousFor: "Actor and performance artist", WikiURL: "https://en.wikipedia.org/wiki/Shia_LaBeouf"},
		{Name: "Hugh Laurie", FamousFor: "Actor ('Dr. Gregory House' in 'House')", WikiURL: "https://en.wikipedia.org/wiki/Hugh_Laurie"},
	},
	"Jun-12": {
		{Name: "George H. W. Bush", FamousFor: "41st U.S. President", WikiURL: "https://en.wikipedia.org/wiki/George_H._W._Bush"},
		{Name: "Anne Frank", FamousFor: "Diarist of 'The Diary of a Young Girl'", WikiURL: "https://en.wikipedia.org/wiki/Anne_Frank"},
		{Name: "Dave Franco", FamousFor: "Actor", WikiURL: "https://en.wikipedia.org/wiki/Dave_Franco"},
	},
	"Jun-13": {
		{Name: "Mary-Kate & Ashley Olsen", FamousFor: "Actresses and fashion designers", WikiURL: "https://en.wikipedia.org/wiki/Mary-Kate_and_Ashley_Olsen"},
		{Name: "Chris Evans", FamousFor: "Actor ('Captain America' in Marvel films)", WikiURL: "https://en.wikipedia.org/wiki/Chris_Evans_(actor)"},
		{Name: "Tim Allen", FamousFor: "Comedian and actor ('Home Improvement')", WikiURL: "https://en.wikipedia.org/wiki/Tim_Allen"},
	},
	"Jun-14": {
		{Name: "Donald Trump", FamousFor: "45th U.S. President and businessman", WikiURL: "https://en.wikipedia.org/wiki/Donald_Trump"},
		{Name: "Boy George", FamousFor: "Lead singer of Culture Club", WikiURL: "https://en.wikipedia.org/wiki/Boy_George"},
		{Name: "Lucy Hale", FamousFor: "Actress ('Pretty Little Liars')", WikiURL: "https://en.wikipedia.org/wiki/Lucy_Hale"},
	},
	"Jun-15": {
		{Name: "Ice Cube", FamousFor: "Rapper and actor", WikiURL: "https://en.wikipedia.org/wiki/Ice_Cube"},
		{Name: "Neil Patrick Harris", FamousFor: "Actor ('How I Met Your Mother')", WikiURL: "https://en.wikipedia.org/wiki/Neil_Patrick_Harris"},
		{Name: "Courteney Cox", FamousFor: "Actress ('Monica Geller' in 'Friends')", WikiURL: "https://en.wikipedia.org/wiki/Courteney_Cox"},
	},
	"Jun-16": {
		{Name: "Tupac Shakur", FamousFor: "Influential rapper and actor", WikiURL: "https://en.wikipedia.org/wiki/Tupac_Shakur"},
		{Name: "Phil Mickelson", FamousFor: "Professional golfer", WikiURL: "https://en.wikipedia.org/wiki/Phil_Mickelson"},
		{Name: "Daniel Brühl", FamousFor: "Actor", WikiURL: "https://en.wikipedia.org/wiki/Daniel_Br%C3%BChl"},
	},
	"Jun-17": {
		{Name: "Kendrick Lamar", FamousFor: "Pulitzer Prize-winning rapper and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Kendrick_Lamar"},
		{Name: "Venus Williams", FamousFor: "World champion professional tennis player", WikiURL: "https://en.wikipedia.org/wiki/Venus_Williams"},
		{Name: "Barry Manilow", FamousFor: "Singer & Songwriter ('Mandy', 'Copacabana')", WikiURL: "https://en.wikipedia.org/wiki/Barry_Manilow"},
	},
	"Jun-18": {
		{Name: "Paul McCartney", FamousFor: "Bassist and singer for The Beatles", WikiURL: "https://en.wikipedia.org/wiki/Paul_McCartney"},
		{Name: "Blake Shelton", FamousFor: "Country music singer and television personality", WikiURL: "https://en.wikipedia.org/wiki/Blake_Shelton"},
		{Name: "Isabella Rossellini", FamousFor: "Actress and model", WikiURL: "https://en.wikipedia.org/wiki/Isabella_Rossellini"},
	},
	"Jun-19": {
		{Name: "Boris Johnson", FamousFor: "Former Prime Minister of the United Kingdom", WikiURL: "https://en.wikipedia.org/wiki/Boris_Johnson"},
		{Name: "Zoe Saldaña", FamousFor: "Actress ('Avatar', 'Guardians of the Galaxy')", WikiURL: "https://en.wikipedia.org/wiki/Zoe_Salda%C3%B1a"},
		{Name: "Jean Dujardin", FamousFor: "Oscar-winning actor ('The Artist')", WikiURL: "https://en.wikipedia.org/wiki/Jean_Dujardin"},
	},
	"Jun-20": {
		{Name: "Nicole Kidman", FamousFor: "Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Nicole_Kidman"},
		{Name: "Lionel Richie", FamousFor: "Grammy-winning singer and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Lionel_Richie"},
		{Name: "John Goodman", FamousFor: "Actor", WikiURL: "https://en.wikipedia.org/wiki/John_Goodman"},
	},
	"Jun-21": {
		{Name: "Prince William, Prince of Wales", FamousFor: "Heir apparent to the British throne", WikiURL: "https://en.wikipedia.org/wiki/William,_Prince_of_Wales"},
		{Name: "Chris Pratt", FamousFor: "Actor ('Guardians of the Galaxy', 'Jurassic World')", WikiURL: "https://en.wikipedia.org/wiki/Chris_Pratt"},
		{Name: "Lana Del Rey", FamousFor: "Singer and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Lana_Del_Rey"},
	},
	"Jun-22": {
		{Name: "Meryl Streep", FamousFor: "Three-time Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Meryl_Streep"},
		{Name: "Cyndi Lauper", FamousFor: "Grammy-winning singer ('Girls Just Want to Have Fun')", WikiURL: "https://en.wikipedia.org/wiki/Cyndi_Lauper"},
		{Name: "Kris Kristofferson", FamousFor: "Country music singer and actor", WikiURL: "https://en.wikipedia.org/wiki/Kris_Kristofferson"},
	},
	"Jun-23": {
		{Name: "Zinedine Zidane", FamousFor: "French professional football manager and former player", WikiURL: "https://en.wikipedia.org/wiki/Zinedine_Zidane"},
		{Name: "Selma Blair", FamousFor: "Actress", WikiURL: "https://en.wikipedia.org/wiki/Selma_Blair"},
		{Name: "Jason Mraz", FamousFor: "Grammy-winning singer-songwriter ('I'm Yours')", WikiURL: "https://en.wikipedia.org/wiki/Jason_Mraz"},
	},
	"Jun-24": {
		{Name: "Lionel Messi", FamousFor: "Professional footballer, one of the greatest of all time", WikiURL: "https://en.wikipedia.org/wiki/Lionel_Messi"},
		{Name: "Mindy Kaling", FamousFor: "Actress, writer, and producer ('The Office', 'The Mindy Project')", WikiURL: "https://en.wikipedia.org/wiki/Mindy_Kaling"},
		{Name: "Mick Fleetwood", FamousFor: "Drummer and co-founder of Fleetwood Mac", WikiURL: "https://en.wikipedia.org/wiki/Mick_Fleetwood"},
	},
	"Jun-25": {
		{Name: "George Orwell", FamousFor: "Author ('Nineteen Eighty-Four', 'Animal Farm')", WikiURL: "https://en.wikipedia.org/wiki/George_Orwell"},
		{Name: "Anthony Bourdain", FamousFor: "Chef, author, and travel documentarian", WikiURL: "https://en.wikipedia.org/wiki/Anthony_Bourdain"},
		{Name: "Carly Simon", FamousFor: "Singer-songwriter ('You're So Vain')", WikiURL: "https://en.wikipedia.org/wiki/Carly_Simon"},
	},
	"Jun-26": {
		{Name: "Ariana Grande", FamousFor: "Grammy-winning singer and actress", WikiURL: "https://en.wikipedia.org/wiki/Ariana_Grande"},
		{Name: "Derek Jeter", FamousFor: "Hall of Fame baseball shortstop", WikiURL: "https://en.wikipedia.org/wiki/Derek_Jeter"},
		{Name: "Aubrey Plaza", FamousFor: "Actress and comedian ('Parks and Recreation')", WikiURL: "https://en.wikipedia.org/wiki/Aubrey_Plaza"},
	},
	"Jun-27": {
		{Name: "Helen Keller", FamousFor: "Author, political activist, and lecturer; deaf and blind", WikiURL: "https://en.wikipedia.org/wiki/Helen_Keller"},
		{Name: "J.J. Abrams", FamousFor: "Filmmaker ('Star Wars', 'Star Trek', 'Lost')", WikiURL: "https://en.wikipedia.org/wiki/J._J._Abrams"},
		{Name: "Khloé Kardashian", FamousFor: "Media personality and socialite", WikiURL: "https://en.wikipedia.org/wiki/Khlo%C3%A9_Kardashian"},
	},
	"Jun-28": {
		{Name: "Elon Musk", FamousFor: "CEO of SpaceX and Tesla; owner of X (Twitter)", WikiURL: "https://en.wikipedia.org/wiki/Elon_Musk"},
		{Name: "Mel Brooks", FamousFor: "Director, writer, and actor known for comedy films", WikiURL: "https://en.wikipedia.org/wiki/Mel_Brooks"},
		{Name: "John Cusack", FamousFor: "Actor", WikiURL: "https://en.wikipedia.org/wiki/John_Cusack"},
	},
	"Jun-29": {
		{Name: "Nicole Scherzinger", FamousFor: "Lead singer of The Pussycat Dolls", WikiURL: "https://en.wikipedia.org/wiki/Nicole_Scherzinger"},
		{Name: "Gary Busey", FamousFor: "Actor", WikiURL: "https://en.wikipedia.org/wiki/Gary_Busey"},
		{Name: "Antoine de Saint-Exupéry", FamousFor: "Author of 'The Little Prince'", WikiURL: "https://en.wikipedia.org/wiki/Antoine_de_Saint-Exup%C3%A9ry"},
	},
	"Jun-30": {
		{Name: "Michael Phelps", FamousFor: "Most decorated Olympian of all time (swimming)", WikiURL: "https://en.wikipedia.org/wiki/Michael_Phelps"},
		{Name: "Mike Tyson", FamousFor: "Former undisputed heavyweight boxing champion", WikiURL: "https://en.wikipedia.org/wiki/Mike_Tyson"},
		{Name: "Lizzy Caplan", FamousFor: "Actress", WikiURL: "https://en.wikipedia.org/wiki/Lizzy_Caplan"},
	},
	"Jul-01": {
		{Name: "Princess Diana", FamousFor: "Princess of Wales and humanitarian", WikiURL: "https://en.wikipedia.org/wiki/Diana,_Princess_of_Wales"},
		{Name: "Pamela Anderson", FamousFor: "Actress and model ('Baywatch')", WikiURL: "https://en.wikipedia.org/wiki/Pamela_Anderson"},
		{Name: "Liv Tyler", FamousFor: "Actress ('Arwen' in 'Lord of the Rings')", WikiURL: "https://en.wikipedia.org/wiki/Liv_Tyler"},
	},
	"Jul-02": {
		{Name: "Margot Robbie", FamousFor: "Actress ('Barbie', 'I, Tonya')", WikiURL: "https://en.wikipedia.org/wiki/Margot_Robbie"},
		{Name: "Lindsay Lohan", FamousFor: "Actress ('Mean Girls', 'The Parent Trap')", WikiURL: "https://en.wikipedia.org/wiki/Lindsay_Lohan"},
		{Name: "Larry David", FamousFor: "Co-creator of 'Seinfeld', creator and star of 'Curb Your Enthusiasm'", WikiURL: "https://en.wikipedia.org/wiki/Larry_David"},
	},
	"Jul-03": {
		{Name: "Tom Cruise", FamousFor: "Actor and producer ('Top Gun', 'Mission: Impossible')", WikiURL: "https://en.wikipedia.org/wiki/Tom_Cruise"},
		{Name: "Franz Kafka", FamousFor: "Author ('The Metamorphosis')", WikiURL: "https://en.wikipedia.org/wiki/Franz_Kafka"},
		{Name: "Olivia Munn", FamousFor: "Actress and model", WikiURL: "https://en.wikipedia.org/wiki/Olivia_Munn"},
	},
	"Jul-04": {
		{Name: "Calvin Coolidge", FamousFor: "30th U.S. President", WikiURL: "https://en.wikipedia.org/wiki/Calvin_Coolidge"},
		{Name: "Post Malone", FamousFor: "Grammy-nominated rapper, singer, and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Post_Malone"},
		{Name: "Malia Obama", FamousFor: "Daughter of former U.S. President Barack Obama", WikiURL: "https://en.wikipedia.org/wiki/Malia_Obama"},
	},
	"Jul-05": {
		{Name: "P. T. Barnum", FamousFor: "Showman and founder of the Barnum & Bailey Circus", WikiURL: "https://en.wikipedia.org/wiki/P._T._Barnum"},
		{Name: "Eva Green", FamousFor: "Actress ('Casino Royale')", WikiURL: "https://en.wikipedia.org/wiki/Eva_Green"},
		{Name: "Edie Falco", FamousFor: "Emmy-winning actress ('The Sopranos', 'Nurse Jackie')", WikiURL: "https://en.wikipedia.org/wiki/Edie_Falco"},
	},
	"Jul-06": {
		{Name: "George W. Bush", FamousFor: "43rd U.S. President", WikiURL: "https://en.wikipedia.org/wiki/George_W._Bush"},
		{Name: "Kevin Hart", FamousFor: "Comedian and actor", WikiURL: "https://en.wikipedia.org/wiki/Kevin_Hart"},
		{Name: "Sylvester Stallone", FamousFor: "Actor ('Rocky', 'Rambo')", WikiURL: "https://en.wikipedia.org/wiki/Sylvester_Stallone"},
	},
	"Jul-07": {
		{Name: "Ringo Starr", FamousFor: "Drummer for The Beatles", WikiURL: "https://en.wikipedia.org/wiki/Ringo_Starr"},
		{Name: "Dalai Lama (14th)", FamousFor: "Spiritual leader of Tibetan Buddhism", WikiURL: "https://en.wikipedia.org/wiki/14th_Dalai_Lama"},
		{Name: "Michelle Kwan", FamousFor: "Olympic champion figure skater", WikiURL: "https://en.wikipedia.org/wiki/Michelle_Kwan"},
	},
	"Jul-08": {
		{Name: "Kevin Bacon", FamousFor: "Actor ('Footloose')", WikiURL: "https://en.wikipedia.org/wiki/Kevin_Bacon"},
		{Name: "Jaden Smith", FamousFor: "Rapper, singer, and actor", WikiURL: "https://en.wikipedia.org/wiki/Jaden_Smith"},
		{Name: "John D. Rockefeller", FamousFor: "Business magnate and philanthropist", WikiURL: "https://en.wikipedia.org/wiki/John_D._Rockefeller"},
	},
	"Jul-09": {
		{Name: "Tom Hanks", FamousFor: "Two-time Oscar-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Tom_Hanks"},
		{Name: "Courtney Love", FamousFor: "Singer for the band Hole", WikiURL: "https://en.wikipedia.org/wiki/Courtney_Love"},
		{Name: "Jack White", FamousFor: "Musician, frontman of The White Stripes", WikiURL: "https://en.wikipedia.org/wiki/Jack_White"},
	},
	"Jul-10": {
		{Name: "Nikola Tesla", FamousFor: "Inventor and electrical engineer (alternating current)", WikiURL: "https://en.wikipedia.org/wiki/Nikola_Tesla"},
		{Name: "Jessica Simpson", FamousFor: "Singer and actress", WikiURL: "https://en.wikipedia.org/wiki/Jessica_Simpson"},
		{Name: "Chiwetel Ejiofor", FamousFor: "Oscar-nominated actor ('12 Years a Slave')", WikiURL: "https://en.wikipedia.org/wiki/Chiwetel_Ejiofor"},
	},
	"Jul-11": {
		{Name: "John Quincy Adams", FamousFor: "6th U.S. President", WikiURL: "https://en.wikipedia.org/wiki/John_Quincy_Adams"},
		{Name: "Giorgio Armani", FamousFor: "Fashion designer", WikiURL: "https://en.wikipedia.org/wiki/Giorgio_Armani"},
		{Name: "Yul Brynner", FamousFor: "Oscar-winning actor ('The King and I')", WikiURL: "https://en.wikipedia.org/wiki/Yul_Brynner"},
	},
	"Jul-12": {
		{Name: "Henry David Thoreau", FamousFor: "Essayist, poet, and philosopher ('Walden')", WikiURL: "https://en.wikipedia.org/wiki/Henry_David_Thoreau"},
		{Name: "Malala Yousafzai", FamousFor: "Nobel Prize laureate and activist for female education", WikiURL: "https://en.wikipedia.org/wiki/Malala_Yousafzai"},
		{Name: "Michelle Rodriguez", FamousFor: "Actress ('The Fast and the Furious' series)", WikiURL: "https://en.wikipedia.org/wiki/Michelle_Rodriguez"},
	},
	"Jul-13": {
		{Name: "Harrison Ford", FamousFor: "Actor ('Indiana Jones', 'Han Solo' in 'Star Wars')", WikiURL: "https://en.wikipedia.org/wiki/Harrison_Ford"},
		{Name: "Patrick Stewart", FamousFor: "Actor ('Captain Picard' in 'Star Trek', 'Professor X' in 'X-Men')", WikiURL: "https://en.wikipedia.org/wiki/Patrick_Stewart"},
		{Name: "Julius Caesar", FamousFor: "Roman general and statesman", WikiURL: "https://en.wikipedia.org/wiki/Julius_Caesar"},
	},
	"Jul-14": {
		{Name: "Gerald Ford", FamousFor: "38th U.S. President", WikiURL: "https://en.wikipedia.org/wiki/Gerald_Ford"},
		{Name: "Conor McGregor", FamousFor: "UFC champion mixed martial artist", WikiURL: "https://en.wikipedia.org/wiki/Conor_McGregor"},
		{Name: "Jane Lynch", FamousFor: "Emmy-winning actress ('Glee')", WikiURL: "https://en.wikipedia.org/wiki/Jane_Lynch"},
	},
	"Jul-15": {
		{Name: "Rembrandt", FamousFor: "Dutch Golden Age painter", WikiURL: "https://en.wikipedia.org/wiki/Rembrandt"},
		{Name: "Forest Whitaker", FamousFor: "Oscar-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Forest_Whitaker"},
		{Name: "Diane Kruger", FamousFor: "Actress", WikiURL: "https://en.wikipedia.org/wiki/Diane_Kruger"},
	},
	"Jul-16": {
		{Name: "Will Ferrell", FamousFor: "Comedian and actor ('Anchorman')", WikiURL: "https://en.wikipedia.org/wiki/Will_Ferrell"},
		{Name: "Roald Amundsen", FamousFor: "Explorer who led the first expedition to the South Pole", WikiURL: "https://en.wikipedia.org/wiki/Roald_Amundsen"},
		{Name: "Phoebe Cates", FamousFor: "Actress ('Fast Times at Ridgemont High')", WikiURL: "https://en.wikipedia.org/wiki/Phoebe_Cates"},
	},
	"Jul-17": {
		{Name: "Angela Merkel", FamousFor: "Former Chancellor of Germany", WikiURL: "https://en.wikipedia.org/wiki/Angela_Merkel"},
		{Name: "David Hasselhoff", FamousFor: "Actor ('Baywatch', 'Knight Rider')", WikiURL: "https://en.wikipedia.org/wiki/David_Hasselhoff"},
		{Name: "Camilla, Queen of the United Kingdom", FamousFor: "Queen consort of the United Kingdom", WikiURL: "https://en.wikipedia.org/wiki/Queen_Camilla"},
	},
	"Jul-18": {
		{Name: "Nelson Mandela", FamousFor: "Anti-apartheid revolutionary and former President of South Africa", WikiURL: "https://en.wikipedia.org/wiki/Nelson_Mandela"},
		{Name: "Vin Diesel", FamousFor: "Actor ('The Fast and the Furious' series)", WikiURL: "https://en.wikipedia.org/wiki/Vin_Diesel"},
		{Name: "Priyanka Chopra", FamousFor: "Actress and producer", WikiURL: "https://en.wikipedia.org/wiki/Priyanka_Chopra"},
	},
	"Jul-19": {
		{Name: "Benedict Cumberbatch", FamousFor: "Actor ('Sherlock', 'Doctor Strange')", WikiURL: "https://en.wikipedia.org/wiki/Benedict_Cumberbatch"},
		{Name: "Brian May", FamousFor: "Lead guitarist of Queen", WikiURL: "https://en.wikipedia.org/wiki/Brian_May"},
		{Name: "Edgar Degas", FamousFor: "Impressionist artist", WikiURL: "https://en.wikipedia.org/wiki/Edgar_Degas"},
	},
	"Jul-20": {
		{Name: "Gisele Bündchen", FamousFor: "Supermodel", WikiURL: "https://en.wikipedia.org/wiki/Gisele_B%C3%BCndchen"},
		{Name: "Carlos Santana", FamousFor: "Grammy-winning guitarist", WikiURL: "https://en.wikipedia.org/wiki/Carlos_Santana"},
		{Name: "Julianne Hough", FamousFor: "Dancer, singer, and actress", WikiURL: "https://en.wikipedia.org/wiki/Julianne_Hough"},
	},
	"Jul-21": {
		{Name: "Robin Williams", FamousFor: "Oscar-winning actor and comedian", WikiURL: "https://en.wikipedia.org/wiki/Robin_Williams"},
		{Name: "Ernest Hemingway", FamousFor: "Nobel Prize-winning author", WikiURL: "https://en.wikipedia.org/wiki/Ernest_Hemingway"},
		{Name: "Josh Hartnett", FamousFor: "Actor", WikiURL: "https://en.wikipedia.org/wiki/Josh_Hartnett"},
	},
	"Jul-22": {
		{Name: "Selena Gomez", FamousFor: "Singer and actress", WikiURL: "https://en.wikipedia.org/wiki/Selena_Gomez"},
		{Name: "Willem Dafoe", FamousFor: "Actor", WikiURL: "https://en.wikipedia.org/wiki/Willem_Dafoe"},
		{Name: "Prince George of Wales", FamousFor: "Member of the British royal family", WikiURL: "https://en.wikipedia.org/wiki/Prince_George_of_Wales"},
	},
	"Jul-23": {
		{Name: "Daniel Radcliffe", FamousFor: "Actor ('Harry Potter')", WikiURL: "https://en.wikipedia.org/wiki/Daniel_Radcliffe"},
		{Name: "Woody Harrelson", FamousFor: "Emmy-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Woody_Harrelson"},
		{Name: "Philip Seymour Hoffman", FamousFor: "Oscar-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Philip_Seymour_Hoffman"},
	},
	"Jul-24": {
		{Name: "Jennifer Lopez", FamousFor: "Singer, actress, and dancer", WikiURL: "https://en.wikipedia.org/wiki/Jennifer_Lopez"},
		{Name: "Amelia Earhart", FamousFor: "Aviation pioneer", WikiURL: "https://en.wikipedia.org/wiki/Amelia_Earhart"},
		{Name: "Alexandre Dumas", FamousFor: "Author ('The Three Musketeers', 'The Count of Monte Cristo')", WikiURL: "https://en.wikipedia.org/wiki/Alexandre_Dumas"},
	},
	"Jul-25": {
		{Name: "Matt LeBlanc", FamousFor: "Actor ('Joey Tribbiani' in 'Friends')", WikiURL: "https://en.wikipedia.org/wiki/Matt_LeBlanc"},
		{Name: "Iman", FamousFor: "Supermodel and entrepreneur", WikiURL: "https://en.wikipedia.org/wiki/Iman_(model)"},
		{Name: "Estelle Getty", FamousFor: "Actress ('Sophia Petrillo' in 'The Golden Girls')", WikiURL: "https://en.wikipedia.org/wiki/Estelle_Getty"},
	},
	"Jul-26": {
		{Name: "Sandra Bullock", FamousFor: "Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Sandra_Bullock"},
		{Name: "Mick Jagger", FamousFor: "Lead singer of The Rolling Stones", WikiURL: "https://en.wikipedia.org/wiki/Mick_Jagger"},
		{Name: "Jason Statham", FamousFor: "Action film actor", WikiURL: "https://en.wikipedia.org/wiki/Jason_Statham"},
	},
	"Jul-27": {
		{Name: "Alex Rodriguez", FamousFor: "Professional baseball player", WikiURL: "https://en.wikipedia.org/wiki/Alex_Rodriguez"},
		{Name: "Taylor Schilling", FamousFor: "Actress ('Orange Is the New Black')", WikiURL: "https://en.wikipedia.org/wiki/Taylor_Schilling"},
		{Name: "Jordan Spieth", FamousFor: "Professional golfer", WikiURL: "https://en.wikipedia.org/wiki/Jordan_Spieth"},
	},
	"Jul-28": {
		{Name: "Jacqueline Kennedy Onassis", FamousFor: "First Lady of the United States", WikiURL: "https://en.wikipedia.org/wiki/Jacqueline_Kennedy_Onassis"},
		{Name: "Soulja Boy", FamousFor: "Rapper ('Crank That')", WikiURL: "https://en.wikipedia.org/wiki/Soulja_Boy"},
		{Name: "Elizabeth Berkley", FamousFor: "Actress ('Saved by the Bell')", WikiURL: "https://en.wikipedia.org/wiki/Elizabeth_Berkley"},
	},
	"Jul-29": {
		{Name: "Benito Mussolini", FamousFor: "Fascist dictator of Italy", WikiURL: "https://en.wikipedia.org/wiki/Benito_Mussolini"},
		{Name: "Josh Radnor", FamousFor: "Actor ('Ted Mosby' in 'How I Met Your Mother')", WikiURL: "https://en.wikipedia.org/wiki/Josh_Radnor"},
		{Name: "Alexis de Tocqueville", FamousFor: "Political scientist and historian ('Democracy in America')", WikiURL: "https://en.wikipedia.org/wiki/Alexis_de_Tocqueville"},
	},
	"Jul-30": {
		{Name: "Arnold Schwarzenegger", FamousFor: "Actor ('The Terminator') and former Governor of California", WikiURL: "https://en.wikipedia.org/wiki/Arnold_Schwarzenegger"},
		{Name: "Henry Ford", FamousFor: "Founder of the Ford Motor Company", WikiURL: "https://en.wikipedia.org/wiki/Henry_Ford"},
		{Name: "Lisa Kudrow", FamousFor: "Actress ('Phoebe Buffay' in 'Friends')", WikiURL: "https://en.wikipedia.org/wiki/Lisa_Kudrow"},
	},
	"Jul-31": {
		{Name: "J. K. Rowling", FamousFor: "Author of the 'Harry Potter' series", WikiURL: "https://en.wikipedia.org/wiki/J._K._Rowling"},
		{Name: "Wesley Snipes", FamousFor: "Actor ('Blade')", WikiURL: "https://en.wikipedia.org/wiki/Wesley_Snipes"},
		{Name: "Mark Cuban", FamousFor: "Billionaire entrepreneur and owner of the Dallas Mavericks", WikiURL: "https://en.wikipedia.org/wiki/Mark_Cuban"},
	},
	"Aug-01": {
		{Name: "Herman Melville", FamousFor: "Author ('Moby-Dick')", WikiURL: "https://en.wikipedia.org/wiki/Herman_Melville"},
		{Name: "Jason Momoa", FamousFor: "Actor ('Aquaman')", WikiURL: "https://en.wikipedia.org/wiki/Jason_Momoa"},
		{Name: "Yves Saint Laurent", FamousFor: "Fashion designer", WikiURL: "https://en.wikipedia.org/wiki/Yves_Saint_Laurent_(designer)"},
	},
	"Aug-02": {
		{Name: "James Baldwin", FamousFor: "Novelist, playwright, and activist", WikiURL: "https://en.wikipedia.org/wiki/James_Baldwin"},
		{Name: "Kevin Smith", FamousFor: "Filmmaker ('Clerks')", WikiURL: "https://en.wikipedia.org/wiki/Kevin_Smith"},
		{Name: "Mary-Louise Parker", FamousFor: "Actress ('Weeds')", WikiURL: "https://en.wikipedia.org/wiki/Mary-Louise_Parker"},
	},
	"Aug-03": {
		{Name: "Tom Brady", FamousFor: "7-time Super Bowl champion NFL quarterback", WikiURL: "https://en.wikipedia.org/wiki/Tom_Brady"},
		{Name: "Martha Stewart", FamousFor: "Businesswoman and television personality", WikiURL: "https://en.wikipedia.org/wiki/Martha_Stewart"},
		{Name: "Martin Sheen", FamousFor: "Actor ('The West Wing')", WikiURL: "https://en.wikipedia.org/wiki/Martin_Sheen"},
	},
	"Aug-04": {
		{Name: "Barack Obama", FamousFor: "44th U.S. President", WikiURL: "https://en.wikipedia.org/wiki/Barack_Obama"},
		{Name: "Meghan, Duchess of Sussex", FamousFor: "Actress and member of the British royal family", WikiURL: "https://en.wikipedia.org/wiki/Meghan,_Duchess_of_Sussex"},
		{Name: "Louis Armstrong", FamousFor: "Influential jazz trumpeter and singer", WikiURL: "https://en.wikipedia.org/wiki/Louis_Armstrong"},
	},
	"Aug-05": {
		{Name: "Neil Armstrong", FamousFor: "Astronaut, first person to walk on the Moon", WikiURL: "https://en.wikipedia.org/wiki/Neil_Armstrong"},
		{Name: "Patrick Ewing", FamousFor: "Hall of Fame basketball player", WikiURL: "https://en.wikipedia.org/wiki/Patrick_Ewing"},
		{Name: "Jesse Williams", FamousFor: "Actor ('Grey's Anatomy')", WikiURL: "https://en.wikipedia.org/wiki/Jesse_Williams"},
	},
	"Aug-06": {
		{Name: "Andy Warhol", FamousFor: "Leading figure in the pop art movement", WikiURL: "https://en.wikipedia.org/wiki/Andy_Warhol"},
		{Name: "Lucille Ball", FamousFor: "Actress and comedian ('I Love Lucy')", WikiURL: "https://en.wikipedia.org/wiki/Lucille_Ball"},
		{Name: "M. Night Shyamalan", FamousFor: "Filmmaker ('The Sixth Sense')", WikiURL: "https://en.wikipedia.org/wiki/M._Night_Shyamalan"},
	},
	"Aug-07": {
		{Name: "Charlize Theron", FamousFor: "Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Charlize_Theron"},
		{Name: "David Duchovny", FamousFor: "Actor ('The X-Files')", WikiURL: "https://en.wikipedia.org/wiki/David_Duchovny"},
		{Name: "Mata Hari", FamousFor: "Exotic dancer convicted of being a spy", WikiURL: "https://en.wikipedia.org/wiki/Mata_Hari"},
	},
	"Aug-08": {
		{Name: "Roger Federer", FamousFor: "Professional tennis player, 20-time Grand Slam winner", WikiURL: "https://en.wikipedia.org/wiki/Roger_Federer"},
		{Name: "Dustin Hoffman", FamousFor: "Two-time Oscar-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Dustin_Hoffman"},
		{Name: "Princess Beatrice", FamousFor: "Member of the British royal family", WikiURL: "https://en.wikipedia.org/wiki/Princess_Beatrice"},
	},
	"Aug-09": {
		{Name: "Whitney Houston", FamousFor: "Grammy-winning singer and actress", WikiURL: "https://en.wikipedia.org/wiki/Whitney_Houston"},
		{Name: "Gillian Anderson", FamousFor: "Actress ('The X-Files', 'The Crown')", WikiURL: "https://en.wikipedia.org/wiki/Gillian_Anderson"},
		{Name: "Anna Kendrick", FamousFor: "Actress and singer ('Pitch Perfect')", WikiURL: "https://en.wikipedia.org/wiki/Anna_Kendrick"},
	},
	"Aug-10": {
		{Name: "Kylie Jenner", FamousFor: "Media personality and businesswoman", WikiURL: "https://en.wikipedia.org/wiki/Kylie_Jenner"},
		{Name: "Antonio Banderas", FamousFor: "Actor", WikiURL: "https://en.wikipedia.org/wiki/Antonio_Banderas"},
		{Name: "Justin Theroux", FamousFor: "Actor and screenwriter", WikiURL: "https://en.wikipedia.org/wiki/Justin_Theroux"},
	},
	"Aug-11": {
		{Name: "Chris Hemsworth", FamousFor: "Actor ('Thor' in Marvel films)", WikiURL: "https://en.wikipedia.org/wiki/Chris_Hemsworth"},
		{Name: "Hulk Hogan", FamousFor: "Professional wrestler and television personality", WikiURL: "https://en.wikipedia.org/wiki/Hulk_Hogan"},
		{Name: "Steve Wozniak", FamousFor: "Co-founder of Apple Inc.", WikiURL: "https://en.wikipedia.org/wiki/Steve_Wozniak"},
	},
	"Aug-12": {
		{Name: "Cara Delevingne", FamousFor: "Model and actress", WikiURL: "https://en.wikipedia.org/wiki/Cara_Delevingne"},
		{Name: "Casey Affleck", FamousFor: "Oscar-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Casey_Affleck"},
		{Name: "George Soros", FamousFor: "Billionaire investor and philanthropist", WikiURL: "https://en.wikipedia.org/wiki/George_Soros"},
	},
	"Aug-13": {
		{Name: "Alfred Hitchcock", FamousFor: "Film director, 'The Master of Suspense'", WikiURL: "https://en.wikipedia.org/wiki/Alfred_Hitchcock"},
		{Name: "Fidel Castro", FamousFor: "Former Prime Minister and President of Cuba", WikiURL: "https://en.wikipedia.org/wiki/Fidel_Castro"},
		{Name: "Sebastian Stan", FamousFor: "Actor ('Bucky Barnes' in Marvel films)", WikiURL: "https://en.wikipedia.org/wiki/Sebastian_Stan"},
	},
	"Aug-14": {
		{Name: "Magic Johnson", FamousFor: "Hall of Fame basketball player", WikiURL: "https://en.wikipedia.org/wiki/Magic_Johnson"},
		{Name: "Halle Berry", FamousFor: "Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Halle_Berry"},
		{Name: "Mila Kunis", FamousFor: "Actress", WikiURL: "https://en.wikipedia.org/wiki/Mila_Kunis"},
	},
	"Aug-15": {
		{Name: "Napoleon Bonaparte", FamousFor: "Emperor of the French", WikiURL: "https://en.wikipedia.org/wiki/Napoleon"},
		{Name: "Jennifer Lawrence", FamousFor: "Oscar-winning actress ('The Hunger Games')", WikiURL: "https://en.wikipedia.org/wiki/Jennifer_Lawrence"},
		{Name: "Ben Affleck", FamousFor: "Oscar-winning actor and director ('Batman')", WikiURL: "https://en.wikipedia.org/wiki/Ben_Affleck"},
	},
	"Aug-16": {
		{Name: "Madonna", FamousFor: "Singer, 'Queen of Pop'", WikiURL: "https://en.wikipedia.org/wiki/Madonna"},
		{Name: "Steve Carell", FamousFor: "Actor and comedian ('Michael Scott' in 'The Office')", WikiURL: "https://en.wikipedia.org/wiki/Steve_Carell"},
		{Name: "James Cameron", FamousFor: "Oscar-winning director ('Titanic', 'Avatar')", WikiURL: "https://en.wikipedia.org/wiki/James_Cameron"},
	},
	"Aug-17": {
		{Name: "Robert De Niro", FamousFor: "Two-time Oscar-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Robert_De_Niro"},
		{Name: "Sean Penn", FamousFor: "Two-time Oscar-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Sean_Penn"},
		{Name: "Davy Crockett", FamousFor: "American folk hero, frontiersman, and politician", WikiURL: "https://en.wikipedia.org/wiki/Davy_Crockett"},
	},
	"Aug-18": {
		{Name: "Robert Redford", FamousFor: "Actor and director", WikiURL: "https://en.wikipedia.org/wiki/Robert_Redford"},
		{Name: "Patrick Swayze", FamousFor: "Actor ('Dirty Dancing', 'Ghost')", WikiURL: "https://en.wikipedia.org/wiki/Patrick_Swayze"},
		{Name: "Andy Samberg", FamousFor: "Comedian and actor ('Brooklyn Nine-Nine')", WikiURL: "https://en.wikipedia.org/wiki/Andy_Samberg"},
	},
	"Aug-19": {
		{Name: "Bill Clinton", FamousFor: "42nd U.S. President", WikiURL: "https://en.wikipedia.org/wiki/Bill_Clinton"},
		{Name: "Matthew Perry", FamousFor: "Actor ('Chandler Bing' in 'Friends')", WikiURL: "https://en.wikipedia.org/wiki/Matthew_Perry"},
		{Name: "Coco Chanel", FamousFor: "Fashion designer", WikiURL: "https://en.wikipedia.org/wiki/Coco_Chanel"},
	},
	"Aug-20": {
		{Name: "Demi Lovato", FamousFor: "Singer and actress", WikiURL: "https://en.wikipedia.org/wiki/Demi_Lovato"},
		{Name: "Andrew Garfield", FamousFor: "Actor ('Spider-Man', 'Hacksaw Ridge')", WikiURL: "https://en.wikipedia.org/wiki/Andrew_Garfield"},
		{Name: "H. P. Lovecraft", FamousFor: "Author of weird and horror fiction", WikiURL: "https://en.wikipedia.org/wiki/H._P._Lovecraft"},
	},
	"Aug-21": {
		{Name: "Usain Bolt", FamousFor: "Olympic champion sprinter, world record holder", WikiURL: "https://en.wikipedia.org/wiki/Usain_Bolt"},
		{Name: "Hayden Panettiere", FamousFor: "Actress ('Heroes')", WikiURL: "https://en.wikipedia.org/wiki/Hayden_Panettiere"},
		{Name: "Kenny Rogers", FamousFor: "Country music singer", WikiURL: "https://en.wikipedia.org/wiki/Kenny_Rogers"},
	},
	"Aug-22": {
		{Name: "Kristen Wiig", FamousFor: "Actress and comedian ('Bridesmaids', 'Saturday Night Live')", WikiURL: "https://en.wikipedia.org/wiki/Kristen_Wiig"},
		{Name: "Dua Lipa", FamousFor: "Grammy-winning singer", WikiURL: "https://en.wikipedia.org/wiki/Dua_Lipa"},
		{Name: "James Corden", FamousFor: "Comedian and television host", WikiURL: "https://en.wikipedia.org/wiki/James_Corden"},
	},
	"Aug-23": {
		{Name: "Kobe Bryant", FamousFor: "Legendary professional basketball player", WikiURL: "https://en.wikipedia.org/wiki/Kobe_Bryant"},
		{Name: "River Phoenix", FamousFor: "Actor", WikiURL: "https://en.wikipedia.org/wiki/River_Phoenix"},
		{Name: "Gene Kelly", FamousFor: "Dancer, actor, and singer ('Singin' in the Rain')", WikiURL: "https://en.wikipedia.org/wiki/Gene_Kelly"},
	},
	"Aug-24": {
		{Name: "Rupert Grint", FamousFor: "Actor ('Ron Weasley' in 'Harry Potter')", WikiURL: "https://en.wikipedia.org/wiki/Rupert_Grint"},
		{Name: "Dave Chappelle", FamousFor: "Comedian", WikiURL: "https://en.wikipedia.org/wiki/Dave_Chappelle"},
		{Name: "Stephen Fry", FamousFor: "Actor, comedian, and writer", WikiURL: "https://en.wikipedia.org/wiki/Stephen_Fry"},
	},
	"Aug-25": {
		{Name: "Sean Connery", FamousFor: "Actor (the first 'James Bond')", WikiURL: "https://en.wikipedia.org/wiki/Sean_Connery"},
		{Name: "Tim Burton", FamousFor: "Filmmaker with a gothic style", WikiURL: "https://en.wikipedia.org/wiki/Tim_Burton"},
		{Name: "Blake Lively", FamousFor: "Actress ('Gossip Girl')", WikiURL: "https://en.wikipedia.org/wiki/Blake_Lively"},
	},
	"Aug-26": {
		{Name: "Macaulay Culkin", FamousFor: "Actor ('Home Alone')", WikiURL: "https://en.wikipedia.org/wiki/Macaulay_Culkin"},
		{Name: "Chris Pine", FamousFor: "Actor ('Star Trek')", WikiURL: "https://en.wikipedia.org/wiki/Chris_Pine"},
		{Name: "Melissa McCarthy", FamousFor: "Actress and comedian", WikiURL: "https://en.wikipedia.org/wiki/Melissa_McCarthy"},
	},
	"Aug-27": {
		{Name: "Aaron Paul", FamousFor: "Emmy-winning actor ('Jesse Pinkman' in 'Breaking Bad')", WikiURL: "https://en.wikipedia.org/wiki/Aaron_Paul"},
		{Name: "Paul Reubens", FamousFor: "Actor and comedian ('Pee-wee Herman')", WikiURL: "https://en.wikipedia.org/wiki/Paul_Reubens"},
		{Name: "Tom Ford", FamousFor: "Fashion designer and filmmaker", WikiURL: "https://en.wikipedia.org/wiki/Tom_Ford"},
	},
	"Aug-28": {
		{Name: "Jack Black", FamousFor: "Actor and musician (Tenacious D)", WikiURL: "https://en.wikipedia.org/wiki/Jack_Black"},
		{Name: "Shania Twain", FamousFor: "Country music singer", WikiURL: "https://en.wikipedia.org/wiki/Shania_Twain"},
		{Name: "Jason Priestley", FamousFor: "Actor ('Beverly Hills, 90210')", WikiURL: "https://en.wikipedia.org/wiki/Jason_Priestley"},
	},
	"Aug-29": {
		{Name: "Michael Jackson", FamousFor: "Singer, 'King of Pop'", WikiURL: "https://en.wikipedia.org/wiki/Michael_Jackson"},
		{Name: "John McCain", FamousFor: "U.S. Senator and presidential candidate", WikiURL: "https://en.wikipedia.org/wiki/John_McCain"},
		{Name: "Lea Michele", FamousFor: "Actress and singer ('Glee')", WikiURL: "https://en.wikipedia.org/wiki/Lea_Michele"},
	},
	"Aug-30": {
		{Name: "Warren Buffett", FamousFor: "Billionaire investor and philanthropist", WikiURL: "https://en.wikipedia.org/wiki/Warren_Buffett"},
		{Name: "Cameron Diaz", FamousFor: "Actress", WikiURL: "https://en.wikipedia.org/wiki/Cameron_Diaz"},
		{Name: "Mary Shelley", FamousFor: "Author ('Frankenstein')", WikiURL: "https://en.wikipedia.org/wiki/Mary_Shelley"},
	},
	"Aug-31": {
		{Name: "Richard Gere", FamousFor: "Actor ('Pretty Woman')", WikiURL: "https://en.wikipedia.org/wiki/Richard_Gere"},
		{Name: "Van Morrison", FamousFor: "Singer-songwriter", WikiURL: "https://en.wikipedia.org/wiki/Van_Morrison"},
		{Name: "Chris Tucker", FamousFor: "Comedian and actor ('Rush Hour')", WikiURL: "https://en.wikipedia.org/wiki/Chris_Tucker"},
	},
	"Sep-01": {
		{Name: "Zendaya", FamousFor: "Emmy-winning actress ('Euphoria', 'Spider-Man')", WikiURL: "https://en.wikipedia.org/wiki/Zendaya"},
		{Name: "Gloria Estefan", FamousFor: "Grammy-winning singer, 'Queen of Latin Pop'", WikiURL: "https://en.wikipedia.org/wiki/Gloria_Estefan"},
		{Name: "Dr. Phil McGraw", FamousFor: "Television host ('Dr. Phil')", WikiURL: "https://en.wikipedia.org/wiki/Phil_McGraw"},
	},
	"Sep-02": {
		{Name: "Keanu Reeves", FamousFor: "Actor ('The Matrix', 'John Wick')", WikiURL: "https://en.wikipedia.org/wiki/Keanu_Reeves"},
		{Name: "Salma Hayek", FamousFor: "Actress and producer", WikiURL: "https://en.wikipedia.org/wiki/Salma_Hayek"},
		{Name: "Mark Harmon", FamousFor: "Actor ('NCIS')", WikiURL: "https://en.wikipedia.org/wiki/Mark_Harmon"},
	},
	"Sep-03": {
		{Name: "Charlie Sheen", FamousFor: "Actor ('Two and a Half Men')", WikiURL: "https://en.wikipedia.org/wiki/Charlie_Sheen"},
		{Name: "Shaun White", FamousFor: "Olympic champion snowboarder and skateboarder", WikiURL: "https://en.wikipedia.org/wiki/Shaun_White"},
		{Name: "Paz de la Huerta", FamousFor: "Actress and model", WikiURL: "https://en.wikipedia.org/wiki/Paz_de_la_Huerta"},
	},
	"Sep-04": {
		{Name: "Beyoncé", FamousFor: "Grammy-winning singer, songwriter, and cultural icon", WikiURL: "https://en.wikipedia.org/wiki/Beyonc%C3%A9"},
		{Name: "Wes Bentley", FamousFor: "Actor ('American Beauty', 'The Hunger Games')", WikiURL: "https://en.wikipedia.org/wiki/Wes_Bentley"},
		{Name: "Ione Skye", FamousFor: "Actress ('Say Anything...')", WikiURL: "https://en.wikipedia.org/wiki/Ione_Skye"},
	},
	"Sep-05": {
		{Name: "Michael Keaton", FamousFor: "Actor ('Batman', 'Birdman')", WikiURL: "https://en.wikipedia.org/wiki/Michael_Keaton"},
		{Name: "Freddie Mercury", FamousFor: "Lead singer of Queen", WikiURL: "https://en.wikipedia.org/wiki/Freddie_Mercury"},
		{Name: "Rose McGowan", FamousFor: "Actress and activist", WikiURL: "https://en.wikipedia.org/wiki/Rose_McGowan"},
	},
	"Sep-06": {
		{Name: "Idris Elba", FamousFor: "Actor ('Luther', 'The Wire')", WikiURL: "https://en.wikipedia.org/wiki/Idris_Elba"},
		{Name: "Jeff Foxworthy", FamousFor: "Comedian ('You might be a redneck if...')", WikiURL: "https://en.wikipedia.org/wiki/Jeff_Foxworthy"},
		{Name: "Roger Waters", FamousFor: "Bassist, co-lead vocalist, and primary lyricist of Pink Floyd", WikiURL: "https://en.wikipedia.org/wiki/Roger_Waters"},
	},
	"Sep-07": {
		{Name: "Queen Elizabeth I", FamousFor: "Queen of England from 1558 to 1603", WikiURL: "https://en.wikipedia.org/wiki/Elizabeth_I"},
		{Name: "Buddy Holly", FamousFor: "Pioneering rock and roll singer-songwriter", WikiURL: "https://en.wikipedia.org/wiki/Buddy_Holly"},
		{Name: "Evan Rachel Wood", FamousFor: "Actress ('Westworld')", WikiURL: "https://en.wikipedia.org/wiki/Evan_Rachel_Wood"},
	},
	"Sep-08": {
		{Name: "Bernie Sanders", FamousFor: "U.S. Senator and presidential candidate", WikiURL: "https://en.wikipedia.org/wiki/Bernie_Sanders"},
		{Name: "Pink", FamousFor: "Grammy-winning singer and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Pink_(singer)"},
		{Name: "David Arquette", FamousFor: "Actor ('Scream' series)", WikiURL: "https://en.wikipedia.org/wiki/David_Arquette"},
	},
	"Sep-09": {
		{Name: "Adam Sandler", FamousFor: "Comedian and actor", WikiURL: "https://en.wikipedia.org/wiki/Adam_Sandler"},
		{Name: "Hugh Grant", FamousFor: "Actor ('Notting Hill', 'Love Actually')", WikiURL: "https://en.wikipedia.org/wiki/Hugh_Grant"},
		{Name: "Michelle Williams", FamousFor: "Oscar-nominated actress", WikiURL: "https://en.wikipedia.org/wiki/Michelle_Williams_(actress)"},
	},
	"Sep-10": {
		{Name: "Colin Firth", FamousFor: "Oscar-winning actor ('The King's Speech')", WikiURL: "https://en.wikipedia.org/wiki/Colin_Firth"},
		{Name: "Arnold Palmer", FamousFor: "Legendary professional golfer", WikiURL: "https://en.wikipedia.org/wiki/Arnold_Palmer"},
		{Name: "Guy Ritchie", FamousFor: "Film director ('Sherlock Holmes', 'Snatch')", WikiURL: "https://en.wikipedia.org/wiki/Guy_Ritchie"},
	},
	"Sep-11": {
		{Name: "Moby", FamousFor: "Musician, DJ, and producer", WikiURL: "https://en.wikipedia.org/wiki/Moby"},
		{Name: "Harry Connick Jr.", FamousFor: "Grammy-winning singer and actor", WikiURL: "https://en.wikipedia.org/wiki/Harry_Connick_Jr."},
		{Name: "Ludacris", FamousFor: "Rapper and actor ('The Fast and the Furious' series)", WikiURL: "https://en.wikipedia.org/wiki/Ludacris"},
	},
	"Sep-12": {
		{Name: "Jesse Owens", FamousFor: "Olympic champion track and field athlete", WikiURL: "https://en.wikipedia.org/wiki/Jesse_Owens"},
		{Name: "Jennifer Hudson", FamousFor: "Oscar-winning actress and Grammy-winning singer", WikiURL: "https://en.wikipedia.org/wiki/Jennifer_Hudson"},
		{Name: "Paul Walker", FamousFor: "Actor ('The Fast and the Furious' series)", WikiURL: "https://en.wikipedia.org/wiki/Paul_Walker"},
	},
	"Sep-13": {
		{Name: "Roald Dahl", FamousFor: "Author ('Charlie and the Chocolate Factory', 'Matilda')", WikiURL: "https://en.wikipedia.org/wiki/Roald_Dahl"},
		{Name: "Fiona Apple", FamousFor: "Grammy-winning singer-songwriter", WikiURL: "https://en.wikipedia.org/wiki/Fiona_Apple"},
		{Name: "Tyler Perry", FamousFor: "Director, actor, and producer ('Madea')", WikiURL: "https://en.wikipedia.org/wiki/Tyler_Perry"},
	},
	"Sep-14": {
		{Name: "Amy Winehouse", FamousFor: "Grammy-winning singer-songwriter", WikiURL: "https://en.wikipedia.org/wiki/Amy_Winehouse"},
		{Name: "Sam Neill", FamousFor: "Actor ('Jurassic Park')", WikiURL: "https://en.wikipedia.org/wiki/Sam_Neill"},
		{Name: "Nas", FamousFor: "Influential rapper", WikiURL: "https://en.wikipedia.org/wiki/Nas"},
	},
	"Sep-15": {
		{Name: "Prince Harry, Duke of Sussex", FamousFor: "Member of the British royal family", WikiURL: "https://en.wikipedia.org/wiki/Prince_Harry,_Duke_of_Sussex"},
		{Name: "Tom Hardy", FamousFor: "Actor ('Venom', 'Mad Max: Fury Road')", WikiURL: "https://en.wikipedia.org/wiki/Tom_Hardy"},
		{Name: "Agatha Christie", FamousFor: "Detective novelist, creator of Hercule Poirot and Miss Marple", WikiURL: "https://en.wikipedia.org/wiki/Agatha_Christie"},
	},
	"Sep-16": {
		{Name: "B.B. King", FamousFor: "Legendary blues guitarist and singer-songwriter", WikiURL: "https://en.wikipedia.org/wiki/B.B._King"},
		{Name: "Mickey Rourke", FamousFor: "Actor ('The Wrestler')", WikiURL: "https://en.wikipedia.org/wiki/Mickey_Rourke"},
		{Name: "Nick Jonas", FamousFor: "Singer, member of the Jonas Brothers", WikiURL: "https://en.wikipedia.org/wiki/Nick_Jonas"},
	},
	"Sep-17": {
		{Name: "Hank Williams", FamousFor: "Pioneering country music singer-songwriter", WikiURL: "https://en.wikipedia.org/wiki/Hank_Williams"},
		{Name: "John Ritter", FamousFor: "Emmy-winning actor ('Three's Company')", WikiURL: "https://en.wikipedia.org/wiki/John_Ritter"},
		{Name: "Kyle Chandler", FamousFor: "Emmy-winning actor ('Friday Night Lights')", WikiURL: "https://en.wikipedia.org/wiki/Kyle_Chandler"},
	},
	"Sep-18": {
		{Name: "Lance Armstrong", FamousFor: "Disgraced professional cyclist", WikiURL: "https://en.wikipedia.org/wiki/Lance_Armstrong"},
		{Name: "James Gandolfini", FamousFor: "Actor ('Tony Soprano' in 'The Sopranos')", WikiURL: "https://en.wikipedia.org/wiki/James_Gandolfini"},
		{Name: "Jada Pinkett Smith", FamousFor: "Actress and talk show host", WikiURL: "https://en.wikipedia.org/wiki/Jada_Pinkett_Smith"},
	},
	"Sep-19": {
		{Name: "Jimmy Fallon", FamousFor: "Comedian and host of 'The Tonight Show'", WikiURL: "https://en.wikipedia.org/wiki/Jimmy_Fallon"},
		{Name: "Jeremy Irons", FamousFor: "Oscar-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Jeremy_Irons"},
		{Name: "Adam West", FamousFor: "Actor ('Batman' in the 1960s TV series)", WikiURL: "https://en.wikipedia.org/wiki/Adam_West"},
	},
	"Sep-20": {
		{Name: "Sophia Loren", FamousFor: "Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Sophia_Loren"},
		{Name: "George R. R. Martin", FamousFor: "Author of 'A Song of Ice and Fire' ('Game of Thrones')", WikiURL: "https://en.wikipedia.org/wiki/George_R._R._Martin"},
		{Name: "Jon Bernthal", FamousFor: "Actor ('The Punisher')", WikiURL: "https://en.wikipedia.org/wiki/Jon_Bernthal"},
	},
	"Sep-21": {
		{Name: "Stephen King", FamousFor: "Author of horror, supernatural fiction, and suspense", WikiURL: "https://en.wikipedia.org/wiki/Stephen_King"},
		{Name: "Bill Murray", FamousFor: "Comedian and actor ('Ghostbusters', 'Lost in Translation')", WikiURL: "https://en.wikipedia.org/wiki/Bill_Murray"},
		{Name: "Jerry Bruckheimer", FamousFor: "Film and television producer", WikiURL: "https://en.wikipedia.org/wiki/Jerry_Bruckheimer"},
	},
	"Sep-22": {
		{Name: "Tom Felton", FamousFor: "Actor ('Draco Malfoy' in 'Harry Potter')", WikiURL: "https://en.wikipedia.org/wiki/Tom_Felton"},
		{Name: "Billie Piper", FamousFor: "Actress and singer ('Doctor Who')", WikiURL: "https://en.wikipedia.org/wiki/Billie_Piper"},
		{Name: "Andrea Bocelli", FamousFor: "Opera singer and tenor", WikiURL: "https://en.wikipedia.org/wiki/Andrea_Bocelli"},
	},
	"Sep-23": {
		{Name: "Bruce Springsteen", FamousFor: "Rock singer-songwriter, 'The Boss'", WikiURL: "https://en.wikipedia.org/wiki/Bruce_Springsteen"},
		{Name: "Ray Charles", FamousFor: "Pioneering soul musician and pianist", WikiURL: "https://en.wikipedia.org/wiki/Ray_Charles"},
		{Name: "Jason Alexander", FamousFor: "Actor ('George Costanza' in 'Seinfeld')", WikiURL: "https://en.wikipedia.org/wiki/Jason_Alexander"},
	},
	"Sep-24": {
		{Name: "F. Scott Fitzgerald", FamousFor: "Author ('The Great Gatsby')", WikiURL: "https://en.wikipedia.org/wiki/F._Scott_Fitzgerald"},
		{Name: "Jim Henson", FamousFor: "Creator of The Muppets", WikiURL: "https://en.wikipedia.org/wiki/Jim_Henson"},
		{Name: "Nia Vardalos", FamousFor: "Writer and star of 'My Big Fat Greek Wedding'", WikiURL: "https://en.wikipedia.org/wiki/Nia_Vardalos"},
	},
	"Sep-25": {
		{Name: "Will Smith", FamousFor: "Grammy and Oscar-winning actor and rapper", WikiURL: "https://en.wikipedia.org/wiki/Will_Smith"},
		{Name: "Michael Douglas", FamousFor: "Oscar-winning actor and producer", WikiURL: "https://en.wikipedia.org/wiki/Michael_Douglas"},
		{Name: "Catherine Zeta-Jones", FamousFor: "Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Catherine_Zeta-Jones"},
	},
	"Sep-26": {
		{Name: "Serena Williams", FamousFor: "Professional tennis player, 23-time Grand Slam winner", WikiURL: "https://en.wikipedia.org/wiki/Serena_Williams"},
		{Name: "T.S. Eliot", FamousFor: "Poet, essayist, and publisher", WikiURL: "https://en.wikipedia.org/wiki/T._S._Eliot"},
		{Name: "Olivia Newton-John", FamousFor: "Singer and actress ('Grease')", WikiURL: "https://en.wikipedia.org/wiki/Olivia_Newton-John"},
	},
	"Sep-27": {
		{Name: "Lil Wayne", FamousFor: "Grammy-winning rapper", WikiURL: "https://en.wikipedia.org/wiki/Lil_Wayne"},
		{Name: "Gwyneth Paltrow", FamousFor: "Oscar-winning actress and founder of Goop", WikiURL: "https://en.wikipedia.org/wiki/Gwyneth_Paltrow"},
		{Name: "Avril Lavigne", FamousFor: "Singer and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Avril_Lavigne"},
	},
	"Sep-28": {
		{Name: "Hilary Duff", FamousFor: "Actress and singer ('Lizzie McGuire')", WikiURL: "https://en.wikipedia.org/wiki/Hilary_Duff"},
		{Name: "Naomi Watts", FamousFor: "Actress", WikiURL: "https://en.wikipedia.org/wiki/Naomi_Watts"},
		{Name: "Confucius", FamousFor: "Chinese philosopher", WikiURL: "https://en.wikipedia.org/wiki/Confucius"},
	},
	"Sep-29": {
		{Name: "Jerry Lee Lewis", FamousFor: "Pioneering rock and roll singer and pianist", WikiURL: "https://en.wikipedia.org/wiki/Jerry_Lee_Lewis"},
		{Name: "Silvio Berlusconi", FamousFor: "Former Prime Minister of Italy", WikiURL: "https://en.wikipedia.org/wiki/Silvio_Berlusconi"},
		{Name: "Kevin Durant", FamousFor: "Championship-winning basketball player", WikiURL: "https://en.wikipedia.org/wiki/Kevin_Durant"},
	},
	"Sep-30": {
		{Name: "Truman Capote", FamousFor: "Author ('In Cold Blood', 'Breakfast at Tiffany's')", WikiURL: "https://en.wikipedia.org/wiki/Truman_Capote"},
		{Name: "Monica Bellucci", FamousFor: "Actress and model", WikiURL: "https://en.wikipedia.org/wiki/Monica_Bellucci"},
		{Name: "Marion Cotillard", FamousFor: "Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Marion_Cotillard"},
	},
	"Oct-01": {
		{Name: "Julie Andrews", FamousFor: "Oscar-winning actress and singer ('Mary Poppins', 'The Sound of Music')", WikiURL: "https://en.wikipedia.org/wiki/Julie_Andrews"},
		{Name: "Jimmy Carter", FamousFor: "39th U.S. President, Nobel Peace Prize laureate", WikiURL: "https://en.wikipedia.org/wiki/Jimmy_Carter"},
		{Name: "Zach Galifianakis", FamousFor: "Comedian and actor ('The Hangover')", WikiURL: "https://en.wikipedia.org/wiki/Zach_Galifianakis"},
	},
	"Oct-02": {
		{Name: "Mahatma Gandhi", FamousFor: "Leader of the Indian independence movement", WikiURL: "https://en.wikipedia.org/wiki/Mahatma_Gandhi"},
		{Name: "Sting", FamousFor: "Grammy-winning musician, lead singer of The Police", WikiURL: "https://en.wikipedia.org/wiki/Sting_(musician)"},
		{Name: "Kelly Ripa", FamousFor: "Television host", WikiURL: "https://en.wikipedia.org/wiki/Kelly_Ripa"},
	},
	"Oct-03": {
		{Name: "Gwen Stefani", FamousFor: "Singer, songwriter, and lead vocalist of No Doubt", WikiURL: "https://en.wikipedia.org/wiki/Gwen_Stefani"},
		{Name: "Alicia Vikander", FamousFor: "Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Alicia_Vikander"},
		{Name: "A$AP Rocky", FamousFor: "Rapper and record producer", WikiURL: "https://en.wikipedia.org/wiki/A$AP_Rocky"},
	},
	"Oct-04": {
		{Name: "Susan Sarandon", FamousFor: "Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Susan_Sarandon"},
		{Name: "Christoph Waltz", FamousFor: "Two-time Oscar-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Christoph_Waltz"},
		{Name: "Dakota Johnson", FamousFor: "Actress ('Fifty Shades' series)", WikiURL: "https://en.wikipedia.org/wiki/Dakota_Johnson"},
	},
	"Oct-05": {
		{Name: "Kate Winslet", FamousFor: "Oscar-winning actress ('Titanic')", WikiURL: "https://en.wikipedia.org/wiki/Kate_Winslet"},
		{Name: "Neil deGrasse Tyson", FamousFor: "Astrophysicist, author, and science communicator", WikiURL: "https://en.wikipedia.org/wiki/Neil_deGrasse_Tyson"},
		{Name: "Jesse Eisenberg", FamousFor: "Actor ('The Social Network')", WikiURL: "https://en.wikipedia.org/wiki/Jesse_Eisenberg"},
	},
	"Oct-06": {
		{Name: "Elisabeth Shue", FamousFor: "Actress ('Leaving Las Vegas', 'The Karate Kid')", WikiURL: "https://en.wikipedia.org/wiki/Elisabeth_Shue"},
		{Name: "Britt Ekland", FamousFor: "Actress and singer", WikiURL: "https://en.wikipedia.org/wiki/Britt_Ekland"},
		{Name: "Amy Jo Johnson", FamousFor: "Actress ('Pink Ranger' in 'Mighty Morphin Power Rangers')", WikiURL: "https://en.wikipedia.org/wiki/Amy_Jo_Johnson"},
	},
	"Oct-07": {
		{Name: "Vladimir Putin", FamousFor: "President of Russia", WikiURL: "https://en.wikipedia.org/wiki/Vladimir_Putin"},
		{Name: "Simon Cowell", FamousFor: "Record executive and television personality ('American Idol', 'The X Factor')", WikiURL: "https://en.wikipedia.org/wiki/Simon_Cowell"},
		{Name: "John Mellencamp", FamousFor: "Rock singer-songwriter", WikiURL: "https://en.wikipedia.org/wiki/John_Mellencamp"},
	},
	"Oct-08": {
		{Name: "Bruno Mars", FamousFor: "Grammy-winning singer, songwriter, and record producer", WikiURL: "https://en.wikipedia.org/wiki/Bruno_Mars"},
		{Name: "Matt Damon", FamousFor: "Oscar-winning actor and screenwriter ('Good Will Hunting')", WikiURL: "https://en.wikipedia.org/wiki/Matt_Damon"},
		{Name: "Sigourney Weaver", FamousFor: "Actress ('Alien', 'Ghostbusters')", WikiURL: "https://en.wikipedia.org/wiki/Sigourney_Weaver"},
	},
	"Oct-09": {
		{Name: "John Lennon", FamousFor: "Singer, songwriter, and member of The Beatles", WikiURL: "https://en.wikipedia.org/wiki/John_Lennon"},
		{Name: "Sharon Osbourne", FamousFor: "Television personality and music manager", WikiURL: "https://en.wikipedia.org/wiki/Sharon_Osbourne"},
		{Name: "Bella Hadid", FamousFor: "Supermodel", WikiURL: "https://en.wikipedia.org/wiki/Bella_Hadid"},
	},
	"Oct-10": {
		{Name: "Brett Favre", FamousFor: "Hall of Fame NFL quarterback", WikiURL: "https://en.wikipedia.org/wiki/Brett_Favre"},
		{Name: "Mario Lopez", FamousFor: "Actor and television host", WikiURL: "https://en.wikipedia.org/wiki/Mario_Lopez"},
		{Name: "Dale Earnhardt Jr.", FamousFor: "Professional auto racing driver", WikiURL: "https://en.wikipedia.org/wiki/Dale_Earnhardt_Jr."},
	},
	"Oct-11": {
		{Name: "Eleanor Roosevelt", FamousFor: "Former First Lady of the United States and diplomat", WikiURL: "https://en.wikipedia.org/wiki/Eleanor_Roosevelt"},
		{Name: "Cardi B", FamousFor: "Grammy-winning rapper", WikiURL: "https://en.wikipedia.org/wiki/Cardi_B"},
		{Name: "Luke Perry", FamousFor: "Actor ('Beverly Hills, 90210')", WikiURL: "https://en.wikipedia.org/wiki/Luke_Perry"},
	},
	"Oct-12": {
		{Name: "Hugh Jackman", FamousFor: "Actor ('Wolverine' in 'X-Men', 'The Greatest Showman')", WikiURL: "https://en.wikipedia.org/wiki/Hugh_Jackman"},
		{Name: "Luciano Pavarotti", FamousFor: "Operatic tenor", WikiURL: "https://en.wikipedia.org/wiki/Luciano_Pavarotti"},
		{Name: "Kirk Cameron", FamousFor: "Actor ('Growing Pains')", WikiURL: "https://en.wikipedia.org/wiki/Kirk_Cameron"},
	},
	"Oct-13": {
		{Name: "Sacha Baron Cohen", FamousFor: "Comedian and actor ('Borat', 'Ali G')", WikiURL: "https://en.wikipedia.org/wiki/Sacha_Baron_Cohen"},
		{Name: "Margaret Thatcher", FamousFor: "Former Prime Minister of the United Kingdom", WikiURL: "https://en.wikipedia.org/wiki/Margaret_Thatcher"},
		{Name: "Paul Simon", FamousFor: "Singer-songwriter (Simon & Garfunkel)", WikiURL: "https://en.wikipedia.org/wiki/Paul_Simon"},
	},
	"Oct-14": {
		{Name: "Usher", FamousFor: "Grammy-winning R&B singer", WikiURL: "https://en.wikipedia.org/wiki/Usher_(musician)"},
		{Name: "Roger Moore", FamousFor: "Actor ('James Bond')", WikiURL: "https://en.wikipedia.org/wiki/Roger_Moore"},
		{Name: "Ralph Lauren", FamousFor: "Fashion designer", WikiURL: "https://en.wikipedia.org/wiki/Ralph_Lauren"},
	},
	"Oct-15": {
		{Name: "Friedrich Nietzsche", FamousFor: "Philosopher", WikiURL: "https://en.wikipedia.org/wiki/Friedrich_Nietzsche"},
		{Name: "Sarah, Duchess of York", FamousFor: "Member of the British royal family", WikiURL: "https://en.wikipedia.org/wiki/Sarah,_Duchess_of_York"},
		{Name: "Emeril Lagasse", FamousFor: "Celebrity chef", WikiURL: "https://en.wikipedia.org/wiki/Emeril_Lagasse"},
	},
	"Oct-16": {
		{Name: "Oscar Wilde", FamousFor: "Poet and playwright ('The Picture of Dorian Gray')", WikiURL: "https://en.wikipedia.org/wiki/Oscar_Wilde"},
		{Name: "Angela Lansbury", FamousFor: "Actress ('Murder, She Wrote')", WikiURL: "https://en.wikipedia.org/wiki/Angela_Lansbury"},
		{Name: "John Mayer", FamousFor: "Grammy-winning singer-songwriter and guitarist", WikiURL: "https://en.wikipedia.org/wiki/John_Mayer"},
	},
	"Oct-17": {
		{Name: "Eminem", FamousFor: "Grammy and Oscar-winning rapper", WikiURL: "https://en.wikipedia.org/wiki/Eminem"},
		{Name: "Rita Hayworth", FamousFor: "Actress and dancer", WikiURL: "https://en.wikipedia.org/wiki/Rita_Hayworth"},
		{Name: "Evel Knievel", FamousFor: "Stunt performer", WikiURL: "https://en.wikipedia.org/wiki/Evel_Knievel"},
	},
	"Oct-18": {
		{Name: "Zac Efron", FamousFor: "Actor ('High School Musical')", WikiURL: "https://en.wikipedia.org/wiki/Zac_Efron"},
		{Name: "Chuck Berry", FamousFor: "Pioneering rock and roll singer and guitarist", WikiURL: "https://en.wikipedia.org/wiki/Chuck_Berry"},
		{Name: "Jean-Claude Van Damme", FamousFor: "Martial artist and actor", WikiURL: "https://en.wikipedia.org/wiki/Jean-Claude_Van_Damme"},
	},
	"Oct-19": {
		{Name: "John le Carré", FamousFor: "Author of espionage novels", WikiURL: "https://en.wikipedia.org/wiki/John_le_Carr%C3%A9"},
		{Name: "Evander Holyfield", FamousFor: "Former undisputed boxing champion", WikiURL: "https://en.wikipedia.org/wiki/Evander_Holyfield"},
		{Name: "Gillian Jacobs", FamousFor: "Actress ('Community')", WikiURL: "https://en.wikipedia.org/wiki/Gillian_Jacobs"},
	},
	"Oct-20": {
		{Name: "Snoop Dogg", FamousFor: "Rapper, songwriter, and media personality", WikiURL: "https://en.wikipedia.org/wiki/Snoop_Dogg"},
		{Name: "Tom Petty", FamousFor: "Rock musician and frontman of Tom Petty and the Heartbreakers", WikiURL: "https://en.wikipedia.org/wiki/Tom_Petty"},
		{Name: "Viggo Mortensen", FamousFor: "Actor ('Aragorn' in 'Lord of the Rings')", WikiURL: "https://en.wikipedia.org/wiki/Viggo_Mortensen"},
	},
	"Oct-21": {
		{Name: "Kim Kardashian", FamousFor: "Media personality, socialite, and businesswoman", WikiURL: "https://en.wikipedia.org/wiki/Kim_Kardashian"},
		{Name: "Carrie Fisher", FamousFor: "Actress ('Princess Leia' in 'Star Wars')", WikiURL: "https://en.wikipedia.org/wiki/Carrie_Fisher"},
		{Name: "Alfred Nobel", FamousFor: "Inventor of dynamite and founder of the Nobel Prize", WikiURL: "https://en.wikipedia.org/wiki/Alfred_Nobel"},
	},
	"Oct-22": {
		{Name: "Jeff Goldblum", FamousFor: "Actor ('Jurassic Park')", WikiURL: "https://en.wikipedia.org/wiki/Jeff_Goldblum"},
		{Name: "Christopher Lloyd", FamousFor: "Actor ('Doc Brown' in 'Back to the Future')", WikiURL: "https://en.wikipedia.org/wiki/Christopher_Lloyd"},
		{Name: "Catherine Deneuve", FamousFor: "Actress", WikiURL: "https://en.wikipedia.org/wiki/Catherine_Deneuve"},
	},
	"Oct-23": {
		{Name: "Pelé", FamousFor: "Brazilian football legend, one of the greatest of all time", WikiURL: "https://en.wikipedia.org/wiki/Pel%C3%A9"},
		{Name: "Ryan Reynolds", FamousFor: "Actor ('Deadpool')", WikiURL: "https://en.wikipedia.org/wiki/Ryan_Reynolds"},
		{Name: "Weird Al Yankovic", FamousFor: "Musician known for his song parodies", WikiURL: "https://en.wikipedia.org/wiki/Weird_Al_Yankovic"},
	},
	"Oct-24": {
		{Name: "Drake", FamousFor: "Grammy-winning rapper and singer", WikiURL: "https://en.wikipedia.org/wiki/Drake_(musician)"},
		{Name: "PewDiePie", FamousFor: "YouTube personality", WikiURL: "https://en.wikipedia.org/wiki/PewDiePie"},
		{Name: "Kevin Kline", FamousFor: "Oscar-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Kevin_Kline"},
	},
	"Oct-25": {
		{Name: "Pablo Picasso", FamousFor: "Influential artist, co-founder of Cubism", WikiURL: "https://en.wikipedia.org/wiki/Pablo_Picasso"},
		{Name: "Katy Perry", FamousFor: "Pop singer and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Katy_Perry"},
		{Name: "Ciara", FamousFor: "Grammy-winning singer", WikiURL: "https://en.wikipedia.org/wiki/Ciara"},
	},
	"Oct-26": {
		{Name: "Hillary Clinton", FamousFor: "Former U.S. Secretary of State, First Lady, and presidential candidate", WikiURL: "https://en.wikipedia.org/wiki/Hillary_Clinton"},
		{Name: "Seth MacFarlane", FamousFor: "Creator of 'Family Guy' and 'American Dad!'", WikiURL: "https://en.wikipedia.org/wiki/Seth_MacFarlane"},
		{Name: "Pat Sajak", FamousFor: "Long-time host of 'Wheel of Fortune'", WikiURL: "https://en.wikipedia.org/wiki/Pat_Sajak"},
	},
	"Oct-27": {
		{Name: "Theodore Roosevelt", FamousFor: "26th U.S. President", WikiURL: "https://en.wikipedia.org/wiki/Theodore_Roosevelt"},
		{Name: "John Cleese", FamousFor: "Comedian, actor, and member of Monty Python", WikiURL: "https://en.wikipedia.org/wiki/John_Cleese"},
		{Name: "Kelly Osbourne", FamousFor: "Television personality and singer", WikiURL: "https://en.wikipedia.org/wiki/Kelly_Osbourne"},
	},
	"Oct-28": {
		{Name: "Bill Gates", FamousFor: "Co-founder of Microsoft and philanthropist", WikiURL: "https://en.wikipedia.org/wiki/Bill_Gates"},
		{Name: "Julia Roberts", FamousFor: "Oscar-winning actress ('Pretty Woman')", WikiURL: "https://en.wikipedia.org/wiki/Julia_Roberts"},
		{Name: "Joaquin Phoenix", FamousFor: "Oscar-winning actor ('Joker')", WikiURL: "https://en.wikipedia.org/wiki/Joaquin_Phoenix"},
	},
	"Oct-29": {
		{Name: "Winona Ryder", FamousFor: "Actress ('Stranger Things')", WikiURL: "https://en.wikipedia.org/wiki/Winona_Ryder"},
		{Name: "Richard Dreyfuss", FamousFor: "Oscar-winning actor ('Jaws')", WikiURL: "https://en.wikipedia.org/wiki/Richard_Dreyfuss"},
		{Name: "Gabrielle Union", FamousFor: "Actress", WikiURL: "https://en.wikipedia.org/wiki/Gabrielle_Union"},
	},
	"Oct-30": {
		{Name: "Diego Maradona", FamousFor: "Argentine football legend", WikiURL: "https://en.wikipedia.org/wiki/Diego_Maradona"},
		{Name: "John Adams", FamousFor: "2nd U.S. President", WikiURL: "https://en.wikipedia.org/wiki/John_Adams"},
		{Name: "Ivanka Trump", FamousFor: "Businesswoman and former presidential advisor", WikiURL: "https://en.wikipedia.org/wiki/Ivanka_Trump"},
	},
	"Oct-31": {
		{Name: "Peter Jackson", FamousFor: "Oscar-winning director ('The Lord of the Rings')", WikiURL: "https://en.wikipedia.org/wiki/Peter_Jackson"},
		{Name: "Vanilla Ice", FamousFor: "Rapper ('Ice Ice Baby')", WikiURL: "https://en.wikipedia.org/wiki/Vanilla_Ice"},
		{Name: "Willow Smith", FamousFor: "Singer and actress", WikiURL: "https://en.wikipedia.org/wiki/Willow_Smith"},
	},
	"Nov-01": {
		{Name: "Tim Cook", FamousFor: "CEO of Apple Inc.", WikiURL: "https://en.wikipedia.org/wiki/Tim_Cook"},
		{Name: "Anthony Kiedis", FamousFor: "Lead singer of Red Hot Chili Peppers", WikiURL: "https://en.wikipedia.org/wiki/Anthony_Kiedis"},
		{Name: "Lyle Lovett", FamousFor: "Grammy-winning country singer-songwriter", WikiURL: "https://en.wikipedia.org/wiki/Lyle_Lovett"},
	},
	"Nov-02": {
		{Name: "Marie Antoinette", FamousFor: "Last Queen of France before the French Revolution", WikiURL: "https://en.wikipedia.org/wiki/Marie_Antoinette"},
		{Name: "David Schwimmer", FamousFor: "Actor ('Ross Geller' in 'Friends')", WikiURL: "https://en.wikipedia.org/wiki/David_Schwimmer"},
		{Name: "Nelly" , FamousFor: "Grammy-winning rapper and singer", WikiURL: "https://en.wikipedia.org/wiki/Nelly"},
	},
	"Nov-03": {
		{Name: "Charles Bronson", FamousFor: "Actor known for 'tough guy' roles", WikiURL: "https://en.wikipedia.org/wiki/Charles_Bronson"},
		{Name: "Dolph Lundgren", FamousFor: "Actor ('Ivan Drago' in 'Rocky IV')", WikiURL: "https://en.wikipedia.org/wiki/Dolph_Lundgren"},
		{Name: "Kendall Jenner", FamousFor: "Model and media personality", WikiURL: "https://en.wikipedia.org/wiki/Kendall_Jenner"},
	},
	"Nov-04": {
		{Name: "Matthew McConaughey", FamousFor: "Oscar-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Matthew_McConaughey"},
		{Name: "Sean 'Diddy' Combs", FamousFor: "Rapper, producer, and record executive", WikiURL: "https://en.wikipedia.org/wiki/Sean_Combs"},
		{Name: "Laura Bush", FamousFor: "Former First Lady of the United States", WikiURL: "https://en.wikipedia.org/wiki/Laura_Bush"},
	},
	"Nov-05": {
		{Name: "Tilda Swinton", FamousFor: "Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Tilda_Swinton"},
		{Name: "Sam Rockwell", FamousFor: "Oscar-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Sam_Rockwell"},
		{Name: "Art Garfunkel", FamousFor: "Singer, member of Simon & Garfunkel", WikiURL: "https://en.wikipedia.org/wiki/Art_Garfunkel"},
	},
	"Nov-06": {
		{Name: "Emma Stone", FamousFor: "Oscar-winning actress ('La La Land')", WikiURL: "https://en.wikipedia.org/wiki/Emma_Stone"},
		{Name: "Sally Field", FamousFor: "Two-time Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Sally_Field"},
		{Name: "Ethan Hawke", FamousFor: "Actor and writer", WikiURL: "https://en.wikipedia.org/wiki/Ethan_Hawke"},
	},
	"Nov-07": {
		{Name: "Marie Curie", FamousFor: "Two-time Nobel Prize-winning physicist and chemist", WikiURL: "https://en.wikipedia.org/wiki/Marie_Curie"},
		{Name: "David Guetta", FamousFor: "DJ and music producer", WikiURL: "https://en.wikipedia.org/wiki/David_Guetta"},
		{Name: "Lorde", FamousFor: "Grammy-winning singer-songwriter", WikiURL: "https://en.wikipedia.org/wiki/Lorde"},
	},
	"Nov-08": {
		{Name: "Gordon Ramsay", FamousFor: "Celebrity chef and television personality", WikiURL: "https://en.wikipedia.org/wiki/Gordon_Ramsay"},
		{Name: "Tara Reid", FamousFor: "Actress ('American Pie')", WikiURL: "https://en.wikipedia.org/wiki/Tara_Reid"},
		{Name: "Alain Delon", FamousFor: "French actor and filmmaker", WikiURL: "https://en.wikipedia.org/wiki/Alain_Delon"},
	},
	"Nov-09": {
		{Name: "Carl Sagan", FamousFor: "Astronomer, cosmologist, and science communicator", WikiURL: "https://en.wikipedia.org/wiki/Carl_Sagan"},
		{Name: "Lou Ferrigno", FamousFor: "Actor and bodybuilder ('The Incredible Hulk')", WikiURL: "https://en.wikipedia.org/wiki/Lou_Ferrigno"},
		{Name: "Hedy Lamarr", FamousFor: "Actress and inventor", WikiURL: "https://en.wikipedia.org/wiki/Hedy_Lamarr"},
	},
	"Nov-10": {
		{Name: "Martin Luther", FamousFor: "Key figure in the Protestant Reformation", WikiURL: "https://en.wikipedia.org/wiki/Martin_Luther"},
		{Name: "Richard Burton", FamousFor: "Actor", WikiURL: "https://en.wikipedia.org/wiki/Richard_Burton"},
		{Name: "Ellen Pompeo", FamousFor: "Actress ('Meredith Grey' in 'Grey's Anatomy')", WikiURL: "https://en.wikipedia.org/wiki/Ellen_Pompeo"},
	},
	"Nov-11": {
		{Name: "Leonardo DiCaprio", FamousFor: "Oscar-winning actor and film producer", WikiURL: "https://en.wikipedia.org/wiki/Leonardo_DiCaprio"},
		{Name: "Demi Moore", FamousFor: "Actress", WikiURL: "https://en.wikipedia.org/wiki/Demi_Moore"},
		{Name: "Kurt Vonnegut", FamousFor: "Author ('Slaughterhouse-Five')", WikiURL: "https://en.wikipedia.org/wiki/Kurt_Vonnegut"},
	},
	"Nov-12": {
		{Name: "Neil Young", FamousFor: "Singer-songwriter and musician", WikiURL: "https://en.wikipedia.org/wiki/Neil_Young"},
		{Name: "Ryan Gosling", FamousFor: "Actor ('La La Land', 'The Notebook')", WikiURL: "https://en.wikipedia.org/wiki/Ryan_Gosling"},
		{Name: "Anne Hathaway", FamousFor: "Oscar-winning actress ('Les Misérables')", WikiURL: "https://en.wikipedia.org/wiki/Anne_Hathaway"},
	},
	"Nov-13": {
		{Name: "Whoopi Goldberg", FamousFor: "EGOT-winning actress, comedian, and television host", WikiURL: "https://en.wikipedia.org/wiki/Whoopi_Goldberg"},
		{Name: "Jimmy Kimmel", FamousFor: "Comedian and late-night talk show host", WikiURL: "https://en.wikipedia.org/wiki/Jimmy_Kimmel"},
		{Name: "Gerard Butler", FamousFor: "Actor ('300')", WikiURL: "https://en.wikipedia.org/wiki/Gerard_Butler"},
	},
	"Nov-14": {
		{Name: "Charles III", FamousFor: "King of the United Kingdom", WikiURL: "https://en.wikipedia.org/wiki/Charles_III"},
		{Name: "Claude Monet", FamousFor: "Founder of French Impressionist painting", WikiURL: "https://en.wikipedia.org/wiki/Claude_Monet"},
		{Name: "Josh Duhamel", FamousFor: "Actor", WikiURL: "https://en.wikipedia.org/wiki/Josh_Duhamel"},
	},
	"Nov-15": {
		{Name: "Georgia O'Keeffe", FamousFor: "Artist, 'Mother of American modernism'", WikiURL: "https://en.wikipedia.org/wiki/Georgia_O%27Keeffe"},
		{Name: "Sam Waterston", FamousFor: "Actor ('Law & Order')", WikiURL: "https://en.wikipedia.org/wiki/Sam_Waterston"},
		{Name: "Chad Kroeger", FamousFor: "Lead singer of Nickelback", WikiURL: "https://en.wikipedia.org/wiki/Chad_Kroeger"},
	},
	"Nov-16": {
		{Name: "Maggie Gyllenhaal", FamousFor: "Actress and director", WikiURL: "https://en.wikipedia.org/wiki/Maggie_Gyllenhaal"},
		{Name: "Pete Davidson", FamousFor: "Comedian and actor ('Saturday Night Live')", WikiURL: "https://en.wikipedia.org/wiki/Pete_Davidson"},
		{Name: "Shigeru Miyamoto", FamousFor: "Video game designer and producer at Nintendo (Mario, Zelda)", WikiURL: "https://en.wikipedia.org/wiki/Shigeru_Miyamoto"},
	},
	"Nov-17": {
		{Name: "Martin Scorsese", FamousFor: "Oscar-winning film director", WikiURL: "https://en.wikipedia.org/wiki/Martin_Scorsese"},
		{Name: "Danny DeVito", FamousFor: "Actor and comedian", WikiURL: "https://en.wikipedia.org/wiki/Danny_DeVito"},
		{Name: "Rachel McAdams", FamousFor: "Actress ('The Notebook')", WikiURL: "https://en.wikipedia.org/wiki/Rachel_McAdams"},
	},
	"Nov-18": {
		{Name: "Mickey Mouse", FamousFor: "Cartoon character and mascot of The Walt Disney Company", WikiURL: "https://en.wikipedia.org/wiki/Mickey_Mouse"},
		{Name: "Owen Wilson", FamousFor: "Actor and comedian", WikiURL: "https://en.wikipedia.org/wiki/Owen_Wilson"},
		{Name: "Chloë Sevigny", FamousFor: "Actress and fashion icon", WikiURL: "https://en.wikipedia.org/wiki/Chlo%C3%AB_Sevigny"},
	},
	"Nov-19": {
		{Name: "Jodie Foster", FamousFor: "Two-time Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Jodie_Foster"},
		{Name: "Adam Driver", FamousFor: "Actor ('Kylo Ren' in 'Star Wars')", WikiURL: "https://en.wikipedia.org/wiki/Adam_Driver"},
		{Name: "Calvin Klein", FamousFor: "Fashion designer", WikiURL: "https://en.wikipedia.org/wiki/Calvin_Klein_(fashion_designer)"},
	},
	"Nov-20": {
		{Name: "Joe Biden", FamousFor: "46th U.S. President", WikiURL: "https://en.wikipedia.org/wiki/Joe_Biden"},
		{Name: "Robert F. Kennedy", FamousFor: "U.S. Attorney General and Senator", WikiURL: "https://en.wikipedia.org/wiki/Robert_F._Kennedy"},
		{Name: "Bo Derek", FamousFor: "Actress and model", WikiURL: "https://en.wikipedia.org/wiki/Bo_Derek"},
	},
	"Nov-21": {
		{Name: "Voltaire", FamousFor: "Enlightenment writer, historian, and philosopher", WikiURL: "https://en.wikipedia.org/wiki/Voltaire"},
		{Name: "Goldie Hawn", FamousFor: "Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Goldie_Hawn"},
		{Name: "Björk", FamousFor: "Singer, songwriter, and actress", WikiURL: "https://en.wikipedia.org/wiki/Bj%C3%B6rk"},
	},
	"Nov-22": {
		{Name: "Scarlett Johansson", FamousFor: "Actress ('Black Widow' in Marvel films)", WikiURL: "https://en.wikipedia.org/wiki/Scarlett_Johansson"},
		{Name: "Mark Ruffalo", FamousFor: "Actor ('Hulk' in Marvel films)", WikiURL: "https://en.wikipedia.org/wiki/Mark_Ruffalo"},
		{Name: "Jamie Lee Curtis", FamousFor: "Oscar-winning actress ('Halloween')", WikiURL: "https://en.wikipedia.org/wiki/Jamie_Lee_Curtis"},
	},
	"Nov-23": {
		{Name: "Miley Cyrus", FamousFor: "Singer and actress", WikiURL: "https://en.wikipedia.org/wiki/Miley_Cyrus"},
		{Name: "Boris Karloff", FamousFor: "Actor ('Frankenstein's monster')", WikiURL: "https://en.wikipedia.org/wiki/Boris_Karloff"},
		{Name: "Franklin Pierce", FamousFor: "14th U.S. President", WikiURL: "https://en.wikipedia.org/wiki/Franklin_Pierce"},
	},
	"Nov-24": {
		{Name: "Charles Darwin", FamousFor: "Naturalist, theory of evolution (publicly debated birthday)", WikiURL: "https://en.wikipedia.org/wiki/Charles_Darwin"},
		{Name: "Katherine Heigl", FamousFor: "Actress ('Grey's Anatomy')", WikiURL: "https://en.wikipedia.org/wiki/Katherine_Heigl"},
		{Name: "Billy Connolly", FamousFor: "Comedian and actor", WikiURL: "https://en.wikipedia.org/wiki/Billy_Connolly"},
	},
	"Nov-25": {
		{Name: "John F. Kennedy Jr.", FamousFor: "Lawyer, journalist, and son of President John F. Kennedy", WikiURL: "https://en.wikipedia.org/wiki/John_F._Kennedy_Jr."},
		{Name: "Amy Grant", FamousFor: "Grammy-winning singer", WikiURL: "https://en.wikipedia.org/wiki/Amy_Grant"},
		{Name: "Christina Applegate", FamousFor: "Emmy-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Christina_Applegate"},
	},
	"Nov-26": {
		{Name: "Tina Turner", FamousFor: "Singer, 'Queen of Rock 'n' Roll'", WikiURL: "https://en.wikipedia.org/wiki/Tina_Turner"},
		{Name: "Charles M. Schulz", FamousFor: "Cartoonist, creator of 'Peanuts'", WikiURL: "https://en.wikipedia.org/wiki/Charles_M._Schulz"},
		{Name: "Rita Ora", FamousFor: "Singer and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Rita_Ora"},
	},
	"Nov-27": {
		{Name: "Bruce Lee", FamousFor: "Martial artist and actor", WikiURL: "https://en.wikipedia.org/wiki/Bruce_Lee"},
		{Name: "Jimi Hendrix", FamousFor: "Pioneering electric guitarist, singer, and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Jimi_Hendrix"},
		{Name: "Bill Nye", FamousFor: "Science educator, 'The Science Guy'", WikiURL: "https://en.wikipedia.org/wiki/Bill_Nye"},
	},
	"Nov-28": {
		{Name: "Jon Stewart", FamousFor: "Comedian and former host of 'The Daily Show'", WikiURL: "https://en.wikipedia.org/wiki/Jon_Stewart"},
		{Name: "Anna Nicole Smith", FamousFor: "Model and television personality", WikiURL: "https://en.wikipedia.org/wiki/Anna_Nicole_Smith"},
		{Name: "Judd Nelson", FamousFor: "Actor ('The Breakfast Club')", WikiURL: "https://en.wikipedia.org/wiki/Judd_Nelson"},
	},
	"Nov-29": {
		{Name: "C. S. Lewis", FamousFor: "Author ('The Chronicles of Narnia')", WikiURL: "https://en.wikipedia.org/wiki/C._S._Lewis"},
		{Name: "Howie Mandel", FamousFor: "Comedian and television host", WikiURL: "https://en.wikipedia.org/wiki/Howie_Mandel"},
		{Name: "Chadwick Boseman", FamousFor: "Actor ('Black Panther')", WikiURL: "https://en.wikipedia.org/wiki/Chadwick_Boseman"},
	},
	"Nov-30": {
		{Name: "Winston Churchill", FamousFor: "Prime Minister of the United Kingdom during WWII", WikiURL: "https://en.wikipedia.org/wiki/Winston_Churchill"},
		{Name: "Mark Twain", FamousFor: "Author ('The Adventures of Tom Sawyer')", WikiURL: "https://en.wikipedia.org/wiki/Mark_Twain"},
		{Name: "Ben Stiller", FamousFor: "Comedian, actor, and director", WikiURL: "https://en.wikipedia.org/wiki/Ben_Stiller"},
	},
	"Dec-01": {
		{Name: "Woody Allen", FamousFor: "Filmmaker, writer, actor, and comedian", WikiURL: "https://en.wikipedia.org/wiki/Woody_Allen"},
		{Name: "Bette Midler", FamousFor: "Singer, songwriter, actress, and comedian", WikiURL: "https://en.wikipedia.org/wiki/Bette_Midler"},
		{Name: "Sarah Silverman", FamousFor: "Comedian, actress, and writer", WikiURL: "https://en.wikipedia.org/wiki/Sarah_Silverman"},
	},
	"Dec-02": {
		{Name: "Britney Spears", FamousFor: "Singer, 'Princess of Pop'", WikiURL: "https://en.wikipedia.org/wiki/Britney_Spears"},
		{Name: "Lucy Liu", FamousFor: "Actress ('Charlie's Angels', 'Kill Bill')", WikiURL: "https://en.wikipedia.org/wiki/Lucy_Liu"},
		{Name: "Aaron Rodgers", FamousFor: "Super Bowl-winning NFL quarterback", WikiURL: "https://en.wikipedia.org/wiki/Aaron_Rodgers"},
	},
	"Dec-03": {
		{Name: "Ozzy Osbourne", FamousFor: "Lead vocalist of Black Sabbath, 'Prince of Darkness'", WikiURL: "https://en.wikipedia.org/wiki/Ozzy_Osbourne"},
		{Name: "Julianne Moore", FamousFor: "Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Julianne_Moore"},
		{Name: "Amanda Seyfried", FamousFor: "Actress ('Mamma Mia!')", WikiURL: "https://en.wikipedia.org/wiki/Amanda_Seyfried"},
	},
	"Dec-04": {
		{Name: "Jay-Z", FamousFor: "Grammy-winning rapper, songwriter, and businessman", WikiURL: "https://en.wikipedia.org/wiki/Jay-Z"},
		{Name: "Tyra Banks", FamousFor: "Supermodel and television personality ('America's Next Top Model')", WikiURL: "https://en.wikipedia.org/wiki/Tyra_Banks"},
		{Name: "Jeff Bridges", FamousFor: "Oscar-winning actor ('The Big Lebowski')", WikiURL: "https://en.wikipedia.org/wiki/Jeff_Bridges"},
	},
	"Dec-05": {
		{Name: "Walt Disney", FamousFor: "Animation pioneer, creator of Disneyland", WikiURL: "https://en.wikipedia.org/wiki/Walt_Disney"},
		{Name: "Frankie Muniz", FamousFor: "Actor ('Malcolm in the Middle'), race car driver", WikiURL: "https://en.wikipedia.org/wiki/Frankie_Muniz"},
		{Name: "Little Richard", FamousFor: "Rock and roll pioneer", WikiURL: "https://en.wikipedia.org/wiki/Little_Richard"},
	},
	"Dec-06": {
		{Name: "Judd Apatow", FamousFor: "Filmmaker, producer, and comedian", WikiURL: "https://en.wikipedia.org/wiki/Judd_Apatow"},
		{Name: "Tom Hulce", FamousFor: "Actor ('Amadeus')", WikiURL: "https://en.wikipedia.org/wiki/Tom_Hulce"},
		{Name: "Peter Buck", FamousFor: "Co-founder and lead guitarist of R.E.M.", WikiURL: "https://en.wikipedia.org/wiki/Peter_Buck"},
	},
	"Dec-07": {
		{Name: "Larry Bird", FamousFor: "Hall of Fame basketball player", WikiURL: "https://en.wikipedia.org/wiki/Larry_Bird"},
		{Name: "Tom Waits", FamousFor: "Singer-songwriter and musician", WikiURL: "https://en.wikipedia.org/wiki/Tom_Waits"},
		{Name: "Ellen Burstyn", FamousFor: "Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Ellen_Burstyn"},
	},
	"Dec-08": {
		{Name: "Jim Morrison", FamousFor: "Lead singer of The Doors", WikiURL: "https://en.wikipedia.org/wiki/Jim_Morrison"},
		{Name: "Nicki Minaj", FamousFor: "Rapper, singer, and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Nicki_Minaj"},
		{Name: "Ian Somerhalder", FamousFor: "Actor ('The Vampire Diaries')", WikiURL: "https://en.wikipedia.org/wiki/Ian_Somerhalder"},
	},
	"Dec-09": {
		{Name: "Judi Dench", FamousFor: "Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Judi_Dench"},
		{Name: "Kirk Douglas", FamousFor: "Actor ('Spartacus')", WikiURL: "https://en.wikipedia.org/wiki/Kirk_Douglas"},
		{Name: "Donny Osmond", FamousFor: "Singer and actor", WikiURL: "https://en.wikipedia.org/wiki/Donny_Osmond"},
	},
	"Dec-10": {
		{Name: "Emily Dickinson", FamousFor: "Poet", WikiURL: "https://en.wikipedia.org/wiki/Emily_Dickinson"},
		{Name: "Kenneth Branagh", FamousFor: "Actor and filmmaker", WikiURL: "https://en.wikipedia.org/wiki/Kenneth_Branagh"},
		{Name: "Raven-Symoné", FamousFor: "Actress ('That's So Raven')", WikiURL: "https://en.wikipedia.org/wiki/Raven-Symon%C3%A9"},
	},
	"Dec-11": {
		{Name: "John Kerry", FamousFor: "U.S. Special Presidential Envoy for Climate and former Secretary of State", WikiURL: "https://en.wikipedia.org/wiki/John_Kerry"},
		{Name: "Hailee Steinfeld", FamousFor: "Oscar-nominated actress and singer", WikiURL: "https://en.wikipedia.org/wiki/Hailee_Steinfeld"},
		{Name: "Rider Strong", FamousFor: "Actor ('Boy Meets World')", WikiURL: "https://en.wikipedia.org/wiki/Rider_Strong"},
	},
	"Dec-12": {
		{Name: "Frank Sinatra", FamousFor: "Singer and actor, 'Ol' Blue Eyes'", WikiURL: "https://en.wikipedia.org/wiki/Frank_Sinatra"},
		{Name: "Bob Barker", FamousFor: "Long-time host of 'The Price Is Right'", WikiURL: "https://en.wikipedia.org/wiki/Bob_Barker"},
		{Name: "Jennifer Connelly", FamousFor: "Oscar-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Jennifer_Connelly"},
	},
	"Dec-13": {
		{Name: "Taylor Swift", FamousFor: "Grammy-winning singer-songwriter and global superstar", WikiURL: "https://en.wikipedia.org/wiki/Taylor_Swift"},
		{Name: "Jamie Foxx", FamousFor: "Oscar-winning actor and Grammy-winning musician", WikiURL: "https://en.wikipedia.org/wiki/Jamie_Foxx"},
		{Name: "Dick Van Dyke", FamousFor: "Emmy-winning actor, comedian, and dancer", WikiURL: "https://en.wikipedia.org/wiki/Dick_Van_Dyke"},
	},
	"Dec-14": {
		{Name: "Nostradamus", FamousFor: "Astrologer and reputed seer", WikiURL: "https://en.wikipedia.org/wiki/Nostradamus"},
		{Name: "Vanessa Hudgens", FamousFor: "Actress and singer ('High School Musical')", WikiURL: "https://en.wikipedia.org/wiki/Vanessa_Hudgens"},
		{Name: "Stan Smith", FamousFor: "Professional tennis player and iconic Adidas shoe namesake", WikiURL: "https://en.wikipedia.org/wiki/Stan_Smith"},
	},
	"Dec-15": {
		{Name: "Don Johnson", FamousFor: "Actor ('Miami Vice')", WikiURL: "https://en.wikipedia.org/wiki/Don_Johnson"},
		{Name: "Adam Brody", FamousFor: "Actor ('The O.C.')", WikiURL: "https://en.wikipedia.org/wiki/Adam_Brody"},
		{Name: "Michelle Dockery", FamousFor: "Actress ('Downton Abbey')", WikiURL: "https://en.wikipedia.org/wiki/Michelle_Dockery"},
	},
	"Dec-16": {
		{Name: "Ludwig van Beethoven", FamousFor: "Composer and pianist", WikiURL: "https://en.wikipedia.org/wiki/Ludwig_van_Beethoven"},
		{Name: "Jane Austen", FamousFor: "Novelist ('Pride and Prejudice')", WikiURL: "https://en.wikipedia.org/wiki/Jane_Austen"},
		{Name: "Theo James", FamousFor: "Actor ('Divergent' series)", WikiURL: "https://en.wikipedia.org/wiki/Theo_James"},
	},
	"Dec-17": {
		{Name: "Pope Francis", FamousFor: "Head of the Catholic Church", WikiURL: "https://en.wikipedia.org/wiki/Pope_Francis"},
		{Name: "Milla Jovovich", FamousFor: "Actress and model ('Resident Evil' series)", WikiURL: "https://en.wikipedia.org/wiki/Milla_Jovovich"},
		{Name: "Eugene Levy", FamousFor: "Emmy-winning actor and comedian ('Schitt's Creek')", WikiURL: "https://en.wikipedia.org/wiki/Eugene_Levy"},
	},
	"Dec-18": {
		{Name: "Steven Spielberg", FamousFor: "Oscar-winning film director, producer, and screenwriter", WikiURL: "https://en.wikipedia.org/wiki/Steven_Spielberg"},
		{Name: "Brad Pitt", FamousFor: "Oscar-winning actor and producer", WikiURL: "https://en.wikipedia.org/wiki/Brad_Pitt"},
		{Name: "Christina Aguilera", FamousFor: "Grammy-winning singer and songwriter", WikiURL: "https://en.wikipedia.org/wiki/Christina_Aguilera"},
	},
	"Dec-19": {
		{Name: "Jake Gyllenhaal", FamousFor: "Actor", WikiURL: "https://en.wikipedia.org/wiki/Jake_Gyllenhaal"},
		{Name: "Alyssa Milano", FamousFor: "Actress ('Charmed', 'Who's the Boss?')", WikiURL: "https://en.wikipedia.org/wiki/Alyssa_Milano"},
		{Name: "Édith Piaf", FamousFor: "Singer, France's national chanteuse", WikiURL: "https://en.wikipedia.org/wiki/%C3%89dith_Piaf"},
	},
	"Dec-20": {
		{Name: "Jonah Hill", FamousFor: "Actor, comedian, and filmmaker", WikiURL: "https://en.wikipedia.org/wiki/Jonah_Hill"},
		{Name: "JoJo", FamousFor: "Singer and actress", WikiURL: "https://en.wikipedia.org/wiki/JoJo_(singer)"},
		{Name: "Ashley Cole", FamousFor: "Professional footballer", WikiURL: "https://en.wikipedia.org/wiki/Ashley_Cole"},
	},
	"Dec-21": {
		{Name: "Samuel L. Jackson", FamousFor: "Actor and producer", WikiURL: "https://en.wikipedia.org/wiki/Samuel_L._Jackson"},
		{Name: "Jane Fonda", FamousFor: "Two-time Oscar-winning actress and activist", WikiURL: "https://en.wikipedia.org/wiki/Jane_Fonda"},
		{Name: "Kiefer Sutherland", FamousFor: "Actor ('Jack Bauer' in '24')", WikiURL: "https://en.wikipedia.org/wiki/Kiefer_Sutherland"},
	},
	"Dec-22": {
		{Name: "Ralph Fiennes", FamousFor: "Actor ('Lord Voldemort' in 'Harry Potter', 'Schindler's List')", WikiURL: "https://en.wikipedia.org/wiki/Ralph_Fiennes"},
		{Name: "Meghan Trainor", FamousFor: "Grammy-winning singer-songwriter", WikiURL: "https://en.wikipedia.org/wiki/Meghan_Trainor"},
		{Name: "Jordin Sparks", FamousFor: "Singer, winner of 'American Idol'", WikiURL: "https://en.wikipedia.org/wiki/Jordin_Sparks"},
	},
	"Dec-23": {
		{Name: "Eddie Vedder", FamousFor: "Lead vocalist of Pearl Jam", WikiURL: "https://en.wikipedia.org/wiki/Eddie_Vedder"},
		{Name: "Carla Bruni", FamousFor: "Singer, songwriter, and former First Lady of France", WikiURL: "https://en.wikipedia.org/wiki/Carla_Bruni"},
		{Name: "Susan Lucci", FamousFor: "Emmy-winning actress ('All My Children')", WikiURL: "https://en.wikipedia.org/wiki/Susan_Lucci"},
	},
	"Dec-24": {
		{Name: "Howard Hughes", FamousFor: "Aviator, engineer, and film director", WikiURL: "https://en.wikipedia.org/wiki/Howard_Hughes"},
		{Name: "Ava Gardner", FamousFor: "Actress", WikiURL: "https://en.wikipedia.org/wiki/Ava_Gardner"},
		{Name: "Ricky Martin", FamousFor: "Grammy-winning singer, 'King of Latin Pop'", WikiURL: "https://en.wikipedia.org/wiki/Ricky_Martin"},
	},
	"Dec-25": {
		{Name: "Isaac Newton", FamousFor: "Physicist and mathematician (Laws of Motion)", WikiURL: "https://en.wikipedia.org/wiki/Isaac_Newton"},
		{Name: "Humphrey Bogart", FamousFor: "Actor ('Casablanca')", WikiURL: "https://en.wikipedia.org/wiki/Humphrey_Bogart"},
		{Name: "Justin Trudeau", FamousFor: "Prime Minister of Canada", WikiURL: "https://en.wikipedia.org/wiki/Justin_Trudeau"},
	},
	"Dec-26": {
		{Name: "Jared Leto", FamousFor: "Oscar-winning actor and musician (Thirty Seconds to Mars)", WikiURL: "https://en.wikipedia.org/wiki/Jared_Leto"},
		{Name: "Kit Harington", FamousFor: "Actor ('Jon Snow' in 'Game of Thrones')", WikiURL: "https://en.wikipedia.org/wiki/Kit_Harington"},
		{Name: "Lars Ulrich", FamousFor: "Drummer and co-founder of Metallica", WikiURL: "https://en.wikipedia.org/wiki/Lars_Ulrich"},
	},
	"Dec-27": {
		{Name: "Marlene Dietrich", FamousFor: "Actress and singer", WikiURL: "https://en.wikipedia.org/wiki/Marlene_Dietrich"},
		{Name: "Gerard Depardieu", FamousFor: "Actor", WikiURL: "https://en.wikipedia.org/wiki/G%C3%A9rard_Depardieu"},
		{Name: "Hayley Williams", FamousFor: "Lead vocalist of Paramore", WikiURL: "https://en.wikipedia.org/wiki/Hayley_Williams"},
	},
	"Dec-28": {
		{Name: "Stan Lee", FamousFor: "Comic book writer, editor, and publisher at Marvel Comics", WikiURL: "https://en.wikipedia.org/wiki/Stan_Lee"},
		{Name: "Denzel Washington", FamousFor: "Two-time Oscar-winning actor", WikiURL: "https://en.wikipedia.org/wiki/Denzel_Washington"},
		{Name: "John Legend", FamousFor: "EGOT-winning singer, songwriter, and record producer", WikiURL: "https://en.wikipedia.org/wiki/John_Legend"},
	},
	"Dec-29": {
		{Name: "Jude Law", FamousFor: "Actor", WikiURL: "https://en.wikipedia.org/wiki/Jude_Law"},
		{Name: "Ted Danson", FamousFor: "Emmy-winning actor ('Cheers', 'The Good Place')", WikiURL: "https://en.wikipedia.org/wiki/Ted_Danson"},
		{Name: "Mary Tyler Moore", FamousFor: "Emmy-winning actress", WikiURL: "https://en.wikipedia.org/wiki/Mary_Tyler_Moore"},
	},
	"Dec-30": {
		{Name: "LeBron James", FamousFor: "Championship-winning professional basketball player", WikiURL: "https://en.wikipedia.org/wiki/LeBron_James"},
		{Name: "Tiger Woods", FamousFor: "Professional golfer", WikiURL: "https://en.wikipedia.org/wiki/Tiger_Woods"},
		{Name: "Tyrese Gibson", FamousFor: "Singer and actor ('The Fast and the Furious' series)", WikiURL: "https://en.wikipedia.org/wiki/Tyrese_Gibson"},
	},
	"Dec-31": {
		{Name: "Henri Matisse", FamousFor: "Artist, a leading figure in modern art", WikiURL: "https://en.wikipedia.org/wiki/Henri_Matisse"},
		{Name: "Anthony Hopkins", FamousFor: "Oscar-winning actor ('The Silence of the Lambs')", WikiURL: "https://en.wikipedia.org/wiki/Anthony_Hopkins"},
		{Name: "Val Kilmer", FamousFor: "Actor ('Top Gun', 'Batman Forever')", WikiURL: "https://en.wikipedia.org/wiki/Val_Kilmer"},
	},

	// Other months can be filled in using the same struct format.
	// For example:
	// "Feb-01": { {Name: "Harry Styles", FamousFor: "Singer", WikiURL:"..."} },
}

// Handler and main functions (no changes)
func greetingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	name := r.FormValue("name")
	birthdayStr := r.FormValue("birthday")
	if name == "" || birthdayStr == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	birthday, err := time.Parse("2006-01-02", birthdayStr)
	if err != nil {
		http.Error(w, "Invalid date format. Please use yy-mm-dd.", http.StatusBadRequest)
		return
	}
	now := time.Now()
	if birthday.After(now) {
		http.Error(w, "Invalid birthday: The date cannot be in the future.", http.StatusBadRequest)
		return
	}
	age := now.Year() - birthday.Year()
	if now.YearDay() < birthday.YearDay() {
		age--
	}
	sundayCount := 0
	var sundayDates []string
	for year := birthday.Year(); year <= now.Year(); year++ {
		birthdayInYear := time.Date(year, birthday.Month(), birthday.Day(), 0, 0, 0, 0, time.UTC)
		if birthdayInYear.After(now) {
			break
		}
		if birthdayInYear.Weekday() == time.Sunday {
			sundayCount++
			formattedDate := birthdayInYear.Format("January 2, 2006")
			sundayDates = append(sundayDates, formattedDate)
		}
	}
	var nextSundayBirthdayStr string
	startYear := now.Year()
	birthdayThisYear := time.Date(startYear, birthday.Month(), birthday.Day(), 0, 0, 0, 0, time.UTC)
	if birthdayThisYear.Before(now) {
		startYear++
	}
	for i := 0; i < 100; i++ {
		yearToSearch := startYear + i
		futureBirthday := time.Date(yearToSearch, birthday.Month(), birthday.Day(), 0, 0, 0, 0, time.UTC)
		if futureBirthday.Weekday() == time.Sunday {
			nextSundayBirthdayStr = futureBirthday.Format("Monday, January 2, 2006")
			break
		}
	}
	lookupKey := birthday.Format("Jan-02")
	celebs := celebrityBirthdays[lookupKey]
	greeting := fmt.Sprintf("Hello, %s!", name)
	if now.Month() == birthday.Month() && now.Day() == birthday.Day() {
		greeting = fmt.Sprintf("Happy Birthday, %s!", name)
	}
	data := PageData{
		Greeting:           greeting,
		Age:                age,
		SundayCount:        sundayCount,
		SundayDates:        sundayDates,
		NextSundayBirthday: nextSundayBirthdayStr,
		Celebrities:        celebs,
	}
	tmpl, err := template.New("greeting").Parse(greetingHTML)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index").Parse(indexHTML)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/greet", greetingHandler)
	fmt.Println("Server starting on http://localhost:9090")
	fmt.Println("Press Ctrl+C to stop the server.")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		panic(err)
	}
}