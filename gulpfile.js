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
const es = require('event-stream');
const debug = process.env.NODE_ENV !== 'production';
let defaultTasks = [ 'browserify', 'sass', 'html', 'img' ] ;

gulp.task('default', defaultTasks);

gulp.task('browserify', function () {
  let files = [
    './assets/src/js/script.js',
    './assets/src/js/tracker.js',
  ];

  var tasks = files.map(function(entry) {
    let stream = browserify({
        entries: entry,
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
    .pipe(source(entry.split('/').pop()))

    if(!debug) {
      stream.pipe(buffer()).pipe(uglify())
    }    

    return stream.pipe(gulp.dest(`./assets/build/js`))  
  });

  // create a merged stream
  return es.merge.apply(null, tasks);
});

gulp.task('img', function() {
  return gulp.src('./assets/src/img/**/*')
    .pipe(gulp.dest(`./assets/build/img`))
});

gulp.task('html', function() {
  return gulp.src('./assets/src/**/*.html')
    .pipe(gulp.dest(`./assets/build/`))
});

gulp.task('sass', function () {
	var files = './assets/src/sass/[^_]*.scss';
	return gulp.src(files)
		.pipe(sass())
    .on('error', gutil.log)
		.pipe(rename({ extname: '.css' }))
		.pipe(gulp.dest(`./assets/build/css`))
});

gulp.task('watch', ['default'], function() {
  gulp.watch(['./assets/src/js/**/*.js'], ['browserify'] );
  gulp.watch(['./assets/src/sass/**/**/*.scss'], ['sass'] );
  gulp.watch(['./assets/src/**/*.html'], ['html'] );
  gulp.watch(['./assets/src/img/**/*'], ['img'] );
});
