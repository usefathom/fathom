'use strict';

// TODO: Minify JS

const babelify = require("babelify");
const browserify = require('browserify');
const gulp = require('gulp');
const source = require('vinyl-source-stream');
const buffer = require('vinyl-buffer');
const rename = require('gulp-rename');
const gutil = require('gulp-util');
const sass = require('gulp-sass');

gulp.task('default', [ 'browserify', 'sass' ] );

gulp.task('browserify', function () {
    return browserify({
            entries: './assets/js/script.js',
            debug: true
        })
        .transform("babelify", {presets: ["es2015", "react"]})
        .on('error', gutil.log)
        .bundle()
        .pipe(source('script.js'))
        .pipe(buffer())
        .pipe(gulp.dest('./static/js/'))
});

gulp.task('sass', function () {
	var files = './assets/sass/[^_]*.scss';
	return gulp.src(files)
		.pipe(sass())
		.pipe(rename({ extname: '.css' }))
		.pipe(gulp.dest('./static/css'))
});

gulp.task('watch', ['default'], function() {
  gulp.watch(['./assets/js/**/*.js'], ['browserify'] );
  gulp.watch(['./assets/sass/*.scss'], ['sass'] );
});
