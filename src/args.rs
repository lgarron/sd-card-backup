use clap::{CommandFactory, Parser};
use clap_complete::generator::generate;
use clap_complete::{Generator, Shell};
use std::io::stdout;
use std::process::exit;

#[derive(Parser, Debug)]
#[command(author, version, about, long_about = None)]
#[clap(name = "sd-card-backup")]
pub struct SDCardBackupArgs {
    #[clap(long)]
    dry_run: bool,

    /// Print completions for the given shell (instead of generating any icons).
    /// These can be loaded/stored permanently (e.g. when using Homebrew), but they can also be sourced directly, e.g.:
    ///
    ///  folderify --completions fish | source # fish
    ///  source <(folderify --completions zsh) # zsh
    #[clap(long, verbatim_doc_comment, id = "SHELL")]
    completions: Option<Shell>,
}

fn completions_for_shell(cmd: &mut clap::Command, generator: impl Generator) {
    generate(generator, cmd, "sd-card-backup", &mut stdout());
}

pub fn get_args() -> SDCardBackupArgs {
    let mut command = SDCardBackupArgs::command();

    let args = SDCardBackupArgs::parse();
    if let Some(shell) = args.completions {
        completions_for_shell(&mut command, shell);
        exit(0);
    }

    args
}
