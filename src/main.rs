use chrono::{Datelike, Local, NaiveDate};
use clap::{Args, Parser, Subcommand, ValueEnum};
use colored::Colorize;
use dirs::data_dir;
use serde::{Deserialize, Serialize};
use std::fs;
use std::path::PathBuf;

const APP_NAME: &str = "Jett";

#[derive(Parser)]
#[command(
    name = "Jett",
    version,
    about = "Jett - a Catppuccin Mocha themed task tracker",
    long_about = None
)]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// Add a new task
    Add(AddArgs),

    /// Edit an existing task
    Edit(EditArgs),

    /// Delete a task by ID
    Delete(DeleteArgs),

    /// List tasks
    List,

    /// Show a summary of your tasks
    Summary,
}

#[derive(Args)]
struct AddArgs {
    /// Title of the task
    title: String,

    /// Start date (YYYY-MM-DD). Defaults to today if not provided.
    #[arg(short, long)]
    start: Option<String>,

    /// Due date (YYYY-MM-DD)
    #[arg(short, long)]
    due: String,

    /// Priority (high, medium, low) - defaults to medium
    #[arg(short, long, value_enum, default_value_t = Priority::Medium)]
    priority: Priority,
}

#[derive(Args)]
struct EditArgs {
    /// ID of the task to edit
    id: u32,

    /// New title
    #[arg(long)]
    title: Option<String>,

    /// New start date (YYYY-MM-DD)
    #[arg(long)]
    start: Option<String>,

    /// New due date (YYYY-MM-DD)
    #[arg(long)]
    due: Option<String>,

    /// New priority
    #[arg(long, value_enum)]
    priority: Option<Priority>,
}

#[derive(Args)]
struct DeleteArgs {
    /// ID of the task to delete
    id: u32,
}

#[derive(ValueEnum, Clone, Debug, Serialize, Deserialize, PartialEq, Eq)]
#[serde(rename_all = "lowercase")]
enum Priority {
    High,
    Medium,
    Low,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
struct Task {
    id: u32,
    title: String,
    start_date: NaiveDate,
    due_date: NaiveDate,
    priority: Priority,
}

#[derive(Debug, Copy, Clone)]
enum Status {
    OnSchedule,
    DueToday,
    Overdue,
}

fn main() {
    let cli = Cli::parse();

    if let Err(e) = run(cli) {
        eprintln!("{} {}", "Error:".red().bold(), e);
        std::process::exit(1);
    }
}

fn run(cli: Cli) -> Result<(), String> {
    let mut tasks = load_tasks()?;

    match cli.command {
        Commands::Add(args) => {
            cmd_add(&mut tasks, args)?;
            save_tasks(&tasks)?;
        }
        Commands::Edit(args) => {
            cmd_edit(&mut tasks, args)?;
            save_tasks(&tasks)?;
        }
        Commands::Delete(args) => {
            cmd_delete(&mut tasks, args)?;
            save_tasks(&tasks)?;
        }
        Commands::List => {
            cmd_list(&tasks)?;
        }
        Commands::Summary => {
            cmd_summary(&tasks)?;
        }
    }

    Ok(())
}

fn data_file_path() -> Result<PathBuf, String> {
    if let Some(base) = data_dir() {
        let dir = base.join("jett");
        fs::create_dir_all(&dir).map_err(|e| format!("Could not create data dir: {e}"))?;
        Ok(dir.join("tasks.json"))
    } else {
        // Fallback to current directory
        Ok(PathBuf::from("jett_tasks.json"))
    }
}

fn load_tasks() -> Result<Vec<Task>, String> {
    let path = data_file_path()?;
    if !path.exists() {
        return Ok(Vec::new());
    }

    let contents = fs::read_to_string(&path)
        .map_err(|e| format!("Failed to read tasks file {}: {e}", path.display()))?;

    if contents.trim().is_empty() {
        return Ok(Vec::new());
    }

    serde_json::from_str(&contents).map_err(|e| format!("Failed to parse tasks JSON: {e}"))
}

fn save_tasks(tasks: &Vec<Task>) -> Result<(), String> {
    let path = data_file_path()?;
    let json = serde_json::to_string_pretty(tasks)
        .map_err(|e| format!("Failed to serialize tasks: {e}"))?;
    fs::write(&path, json)
        .map_err(|e| format!("Failed to write tasks file {}: {e}", path.display()))
}

fn parse_date(input: &str) -> Result<NaiveDate, String> {
    NaiveDate::parse_from_str(input, "%Y-%m-%d")
        .map_err(|_| format!("Invalid date '{}'. Use YYYY-MM-DD.", input))
}

fn next_id(tasks: &[Task]) -> u32 {
    tasks.iter().map(|t| t.id).max().unwrap_or(0) + 1
}

// ---------- Commands ----------

fn cmd_add(tasks: &mut Vec<Task>, args: AddArgs) -> Result<(), String> {
    let today = Local::now().date_naive();
    let start_date = match args.start {
        Some(s) => parse_date(&s)?,
        None => today,
    };
    let due_date = parse_date(&args.due)?;

    if due_date < start_date {
        return Err("Due date cannot be before start date.".to_string());
    }

    let id = next_id(tasks);

    let task = Task {
        id,
        title: args.title,
        start_date,
        due_date,
        priority: args.priority,
    };

    tasks.push(task);

    println!(
        "{} {}",
        "Added task".truecolor(245, 194, 231).bold(),
        format!("#{}", id).truecolor(245, 194, 231)
    );

    Ok(())
}

fn cmd_edit(tasks: &mut Vec<Task>, args: EditArgs) -> Result<(), String> {
    let task = tasks
        .iter_mut()
        .find(|t| t.id == args.id)
        .ok_or_else(|| format!("Task with ID {} not found.", args.id))?;

    if let Some(title) = args.title {
        task.title = title;
    }

    if let Some(start) = args.start {
        task.start_date = parse_date(&start)?;
    }

    if let Some(due) = args.due {
        task.due_date = parse_date(&due)?;
    }

    if task.due_date < task.start_date {
        return Err("After editing, due date cannot be before start date.".to_string());
    }

    if let Some(priority) = args.priority {
        task.priority = priority;
    }

    println!(
        "{} {}",
        "Updated task".truecolor(245, 194, 231).bold(),
        format!("#{}", task.id).truecolor(245, 194, 231)
    );

    Ok(())
}

fn cmd_delete(tasks: &mut Vec<Task>, args: DeleteArgs) -> Result<(), String> {
    let before = tasks.len();
    tasks.retain(|t| t.id != args.id);
    let after = tasks.len();

    if after == before {
        return Err(format!("Task with ID {} not found.", args.id));
    }

    println!(
        "{} {}",
        "Deleted task".truecolor(245, 194, 231).bold(),
        format!("#{}", args.id).truecolor(245, 194, 231)
    );

    Ok(())
}

fn cmd_list(tasks: &[Task]) -> Result<(), String> {
    if tasks.is_empty() {
        println!(
            "{}",
            "No tasks yet. Add one with `jett add`!".truecolor(245, 194, 231)
        );
        return Ok(());
    }

    let today = Local::now().date_naive();

    println!(
        "{}",
        format!("{} - Task List", APP_NAME)
            .truecolor(245, 194, 231)
            .bold()
    );

    for task in tasks {
        let status = status_for(task, today);
        let status_str = colored_status_label(status);
        let priority_str = colored_priority_label(&task.priority);

        let header = format!("[#{}] {}", task.id, task.title)
            .truecolor(245, 194, 231)
            .bold();

        println!("{header}");
        println!(
            "  {} {}  {} {}  {} {}",
            "Start:".dimmed(),
            format_date(task.start_date).dimmed(),
            "Due:".dimmed(),
            format_date(task.due_date).dimmed(),
            "Priority:".dimmed(),
            priority_str
        );
        println!("  {} {}", "Status:".dimmed(), status_str);
        println!();
    }

    Ok(())
}

fn cmd_summary(tasks: &[Task]) -> Result<(), String> {
    let today = Local::now().date_naive();

    let mut total = 0;
    let mut overdue = 0;
    let mut due_today = 0;
    let mut on_schedule = 0;

    for task in tasks {
        total += 1;
        match status_for(task, today) {
            Status::Overdue => overdue += 1,
            Status::DueToday => due_today += 1,
            Status::OnSchedule => on_schedule += 1,
        }
    }

    println!(
        "{}",
        format!("{} - Summary", APP_NAME)
            .truecolor(245, 194, 231)
            .bold()
    );
    println!(
        "  {} {}",
        "Total tasks:".dimmed(),
        total.to_string().truecolor(245, 194, 231)
    );
    println!(
        "  {} {}",
        "Overdue:".dimmed(),
        overdue.to_string().red().bold()
    );
    println!(
        "  {} {}",
        "Due today:".dimmed(),
        due_today.to_string().yellow().bold()
    );
    println!(
        "  {} {}",
        "On schedule:".dimmed(),
        on_schedule.to_string().green().bold()
    );

    Ok(())
}

// ---------- Status & Colors ----------

fn status_for(task: &Task, today: NaiveDate) -> Status {
    if task.due_date < today {
        Status::Overdue
    } else if task.due_date == today {
        Status::DueToday
    } else {
        Status::OnSchedule
    }
}

fn colored_status_label(status: Status) -> String {
    match status {
        Status::OnSchedule => "On schedule".green().bold().to_string(),
        Status::DueToday => "Due today".yellow().bold().to_string(),
        Status::Overdue => "Behind (overdue)".red().bold().to_string(),
    }
}

fn colored_priority_label(priority: &Priority) -> String {
    match priority {
        Priority::High => "HIGH".red().bold().to_string(),
        Priority::Medium => "MEDIUM".truecolor(249, 226, 175).bold().to_string(), // Catppuccin Mocha yellow-ish
        Priority::Low => "LOW".truecolor(148, 226, 213).bold().to_string(),       // teal-ish
    }
}

fn format_date(date: NaiveDate) -> String {
    format!("{:04}-{:02}-{:02}", date.year(), date.month(), date.day())
}
