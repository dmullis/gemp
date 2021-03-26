#! /bin/sh

set -ue

failing_gofmt=$(git-ls-files '[a-z]*.go' | xargs gofmt -l)
if ! [ -z "${failing_gofmt}" ]
then
    printf "Files failing gofmt:\n%s\n\n" "${failing_gofmt}"
fi

go install

echo Debug:
echo '    dlv debug '$PWD' -- GEMP_ARGS ...'

(
    cat doc/README.1.md

    echo '## gen'
    cat doc/gen-example.md

    echo '## dump'
    cat doc/dump-example.md

    cat doc/epilogue.md
) >README.md

# XX  This depends on flag.PrintDefaults(), which freely commingles space characters
#     and tabs, with eight columns assumed.
#     Git complains about this with fat red rectangles on output of git-diff, git-log, ...
# X   Relative links to these files from within README.md are munged by GitHub
#     during upload:
#       https://docs.github.com/en/github/writing-on-github/basic-writing-and-formatting-syntax#relative-links
gemp -helpAsMarkdown -h     2>doc/usage.md
gemp -helpAsMarkdown -h gen 2>doc/gen-usage.md

# Alternative Markdown processors:
#    1.  'blackfriday'
#    2.  https://pkg.go.dev/github.com/shurcooL/github_flavored_markdown
#        https://github.com/shurcooL/github_flavored_markdown/issues
#    3.  https://docs.github.com/en/rest/reference/markdown
for mdFile in README.md doc/usage.md doc/gen-usage.md
do
    #  --gfm => "GitHub-Flavored-Markdown"
    #    XXX  Despite --gfm, does NOT transform link references to ".md" files into ".html", as
    #         the GitHub website does. 
    marked --gfm $mdFile >${mdFile%%.*}.html
done

echo firefox --new-window README.html
