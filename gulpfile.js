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
      entries: './assets/js/script.js',
      debug: debug
  })
  .transform("babelify", {
    presets: ["es2015"],
    plugins: [ 
      "transform-decorators-legacy", 
      ["transform-react-jsx", { "pragma":"h" } ] 
    ]
  })
  .bundle()
  .on('error', function(err){
    console.log(err.message);
    this.emit('end');
  })
  .pipe(source('script.js'))
  .pipe(buffer())
  .pipe(gulp.dest('./build/js/'))
});

gulp.task('minify', function(cb) {
  process.env.NODE_ENV = 'production';

  pump([
    gulp.src('./build/js/*.js'),
    uglify().on('error', gutil.log),
    gulp.dest('./build/js/')
  ], cb)
});

gulp.task('img', function() {
  return gulp.src('./assets/img/**/*')
    .pipe(gulp.dest('./build/img'))
});

gulp.task('html', function() {
  return gulp.src('./assets/**/*.html')
    .pipe(gulp.dest('./build'))
});

gulp.task('tracker', function() {
  return gulp.src('./assets/js/tracker.js')
    .pipe(gulp.dest('./build/js'))
});

gulp.task('sass', function () {
	var files = './assets/sass/[^_]*.scss';
	return gulp.src(files)
		.pipe(sass())
    .on('error', gutil.log)
		.pipe(rename({ extname: '.css' }))
		.pipe(gulp.dest('./build/css'))
});

gulp.task('watch', ['default'], function() {
  gulp.watch(['./assets/js/**/*.js'], ['browserify', 'tracker'] );
  gulp.watch(['./assets/sass/**/**/*.scss'], ['sass'] );
  gulp.watch(['./assets/**/*.html'], ['html'] );
  gulp.watch(['./assets/img/**/*'], ['img'] );
});
