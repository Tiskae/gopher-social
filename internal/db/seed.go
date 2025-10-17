// Package db: for seeding the DB
package db

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"strconv"

	"github.com/tiskae/go-social/internal/store"
)

var usernamesData = []string{
	"skywander", "techguru", "coffeeaddict", "mountainhiker", "booklover", "codewizard", "nightowl", "digitalnomad", "pixelartist", "urbanexplorer", "stormchaser", "gamerzone", "cosmicdreamer", "dataminer",
	"stargaz ", "quietreader", "cryptonerd", "fastrunner", "cloudsurfer", "logicmaster", "cybersamurai", "forestwalker", "designninja", "retroplayer", "soundseeker", "ironclimber", "happycoder", "wanderlust", "bytehunter", "moonwalker", "aqualover",
	"fitnessjunkie", "solarflare", "swiftthinker", "earthtrekker", "dreamcatcher", "zenmaster", "novastar", "wordsmith", "matrixcoder", "blueocean", "ideaspark", "futurepilot", "chesschampion", "riverunner", "logicbuilder", "spacetraveler", "mindbender", "trailblazer",
}

var titlesData = []string{
	"The Future of Remote Work",
	"Why Morning Routines Matter",
	"Lessons I Learned from Failure",
	"Exploring the Power of Minimalism",
	"How to Stay Productive When Tired",
	"My Journey into Coding",
	"The Art of Saying No",
	"Why Travel Changes Your Perspective",
	"Building Better Habits in 30 Days",
	"The Rise of AI in Everyday Life",
	"What Nature Can Teach Us About Balance",
	"Tips for Mastering Public Speaking",
	"Why Books Still Matter in a Digital Age",
	"How I Overcame Procrastination",
	"The Science of Happiness",
	"Why Creativity Is a Superpower",
	"Lessons from Starting a Side Hustle",
	"The Value of Deep Work",
	"How to Manage Stress Effectively",
	"Why Sleep Is the Real Productivity Hack",
	"My Favorite Tools for Learning Online",
	"Why Curiosity Fuels Growth",
	"The Joy of Lifelong Learning",
	"How Failure Builds Resilience",
	"Why Simplicity Wins in Design",
	"The Importance of Gratitude",
	"How Technology Shapes Our Thinking",
	"The Secret to Building Discipline",
	"Why Reflection Matters for Growth",
	"My Thoughts on the Future of Education",
}

var contentsData = []string{
	"Just wrapped up a super productive day of coding!",
	"Coffee really is the ultimate bug fixer ☕️",
	"Started learning Go today and it feels amazing!",
	"Small wins stack up into big victories.",
	"Sometimes rest is the most productive thing you can do.",
	"Just deployed my first app — what a feeling! 🚀",
	"Books + quiet time = perfect evening.",
	"Traveling opens your mind in ways nothing else can.",
	"Finally mastered Docker after struggling for weeks!",
	"Gratitude turns little into enough.",
	"The hardest part is starting, but it gets easier after that.",
	"Debugging teaches you patience like nothing else.",
	"Nature has a way of resetting your thoughts.",
	"One new skill every month — that's the goal!",
	"Minimalism isn't about less, it's about more clarity.",
	"Sleep early, wake early, feel unstoppable.",
	"Creativity is just connecting old ideas in new ways.",
	"Every failure is just data for the next attempt.",
	"Slow progress is still progress — don't quit.",
	"Today I chose growth over comfort.",
}

var tagsData = []string{
	"coding", "golang", "backend", "api", "database", "sql",
	"productivity", "habits", "focus", "mindset", "growth", "success",
	"travel", "adventure", "lifestyle", "culture", "exploration", "photography",
	"ai", "technology", "future", "innovation", "automation", "machine-learning",
	"reading", "books", "learning", "knowledge", "education", "curiosity",
	"minimalism", "clarity", "design", "simplicity", "organization", "creativity",
	"career", "growth", "skills", "jobs", "opportunities", "development",
	"docker", "devops", "containers", "kubernetes", "cloud", "infrastructure",
	"health", "fitness", "balance", "wellness", "meditation", "nutrition",
}

var commentsData = []string{
	"This is awesome!, Thanks for sharing 🙌",
	"I totally agree with you., Great perspective!",
	"This made my day 😃, I needed to hear this today.",
	"So true, well said., Love this 👏",
	"Interesting take, never thought of it that way., Keep it up, inspiring stuff!",
	"Wow, this hit hard 💯, Can you share more about this?",
	"Really useful insight, thanks!, This resonates with me a lot.",
	"Solid advice, thank you., I've been thinking about this too.",
	"Such a refreshing perspective., Great reminder 🙏",
	"I'm definitely saving this., This is gold 🔥",
	"Nicely explained!, Thanks for putting this into words.",
	"This gave me a new idea., Very well written.",
	"I completely relate to this., Such a valuable post.",
	"Appreciate you sharing this., This is so motivating!",
	"Wow, really well put together., Simple but powerful.",
}

const (
	UsersCount    = 100
	PostsCount    = 200
	CommentsCount = 500
)

func Seed(store store.Storage, db *sql.DB) error {
	ctx := context.Background()

	users := generateUsers(UsersCount)
	tx, _ := db.BeginTx(ctx, nil)
	for i := range UsersCount {
		if err := store.Users.Create(ctx, tx, users[i]); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	tx.Commit()

	posts := generatePosts(200)
	for i := range PostsCount {
		if err := store.Posts.Create(ctx, posts[i]); err != nil {
			return err
		}
	}

	comments := generateComments(500)
	for i := range CommentsCount {
		if err := store.Comments.Create(ctx, comments[i]); err != nil {
			return err
		}
	}

	log.Println("Seeding complete")

	return nil
}

func generateUsers(num int) []*store.User {
	var users []*store.User

	for i := range num {
		rndIdx := rand.Intn(len(usernamesData))

		newUser := &store.User{
			Username: usernamesData[rndIdx] + strconv.Itoa(i),
			Email:    usernamesData[rndIdx] + strconv.Itoa(i) + "@example.com",
			Role:     store.Role{Name: "user"},
		}

		newUser.Password.Set("123123")

		users = append(users, newUser)
	}

	return users
}

func generatePosts(num int) []*store.Post {
	var posts []*store.Post

	for i := range num {
		newPost := &store.Post{
			Title:   titlesData[rand.Intn(len(titlesData))] + " " + strconv.Itoa(i),
			Content: contentsData[rand.Intn(len(contentsData))],
			Tags:    []string{tagsData[rand.Intn(len(tagsData))]},
			UserID:  rand.Int63n(UsersCount) + 1,
		}

		posts = append(posts, newPost)
	}

	return posts
}

func generateComments(num int) []*store.Comment {
	var comments []*store.Comment

	for i := range num {
		newComment := &store.Comment{
			PostID:  rand.Int63n(PostsCount) + 1,
			UserID:  rand.Int63n(UsersCount) + 1,
			Content: commentsData[rand.Intn(len(commentsData))] + " " + strconv.Itoa(i),
		}

		comments = append(comments, newComment)
	}

	return comments
}
