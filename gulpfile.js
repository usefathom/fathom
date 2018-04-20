'use strict';

const babelify = require("babelify")
const browserify = require('browserify')
const gulp = require('gulp')
const source = require('vinyl-source-stream')
const buffer = require('vinyl-buffer')
const rename = require('gulp-rename')
const gutil = require('gulp-util')
const sass = require('gulp-sass')
const uglify = require('gulp-uglify')
const pump = require('pump')
const debug = process.env.NODE_ENV !== 'production';

let defaultTasks = [ 'browserify', 'sass', 'tracker', 'html', 'img' ] ;
if( ! debug ) {
  defaultTasks.push( 'minify' );
}

gulp.task('default', defaultTasks);

gulp.task('browserify', function () {
    return browserify({
            entries: './assets/src/js/script.js',
            debug: debug
        })
        .transform("babelify", {
          presets: ["es2015"],
          plugins: [ ["transform-react-jsx", { "pragma":"h" } ] ]
        })
        .bundle()
        .on('error', function(err){
          console.log(err.message);
          this.emit('end');
        })
        .pipe(source('script.js'))
        .pipe(buffer())
        .pipe(gulp.dest('./assets/dist/js/'))
});

gulp.task('minify', function(cb) {
  process.env.NODE_ENV = 'production';

  pump([
    gulp.src('./assets/dist/js/*.js'),
    uglify().on('error', gutil.log),
    gulp.dest('./assets/dist/js/')
  ], cb)
});

gulp.task('img', function() {
  return gulp.src('./assets/src/img/**/*')
    .pipe(gulp.dest('./assets/dist/img'))
});

gulp.task('html', function() {
  return gulp.src('./assets/src/**/*.html')
    .pipe(gulp.dest('./assets/dist'))
});

gulp.task('tracker', function() {
  return gulp.src('./assets/src/js/tracker.js')
    .pipe(gulp.dest('./assets/dist/js'))
});

gulp.task('sass', function () {
	var files = './assets/src/sass/[^_]*.scss';
	return gulp.src(files)
		.pipe(sass())
    .on('error', gutil.log)
		.pipe(rename({ extname: '.css' }))
		.pipe(gulp.dest('./assets/dist/css'))
});

gulp.task('watch', ['default'], function() {
  gulp.watch(['./assets/src/js/**/*.js'], ['browserify', 'tracker'] );
  gulp.watch(['./assets/src/sass/**/**/*.scss'], ['sass'] );
});
